package dtx

import (
	"errors"
	"log/slog"

	"github.com/mikeschinkel/go-dt"
)

type writer interface {
	Printf(string, ...any)
	Errorf(string, ...any)
}

type MatchBehavior int

const (
	CollectOnMatch         MatchBehavior = 0
	WriteOnMatch           MatchBehavior = 1
	WriteAndCollectOnMatch MatchBehavior = 2
)

type SkipEntryFunc func(dt.DirPath, *dt.DirEntry) bool
type ParseEntryFunc func(dt.EntryPath, dt.DirPath, dt.DirEntry) dt.EntryPath

type DirPathScanner struct {
	DirPath                dt.DirPath
	ContinueOnErr          bool
	IncludeNonRegularFiles bool
	Logger                 *slog.Logger
	Writer                 writer
	MatchBehavior          MatchBehavior
	SkipPaths              []dt.PathSegment
	SkipDirFunc            SkipEntryFunc
	SkipEntryFunc          SkipEntryFunc
	ParseEntryFunc         ParseEntryFunc
}

type DirPathScannerArgs struct {
	DirPath                dt.DirPath
	ContinueOnErr          bool
	IncludeNonRegularFiles bool
	MatchBehavior          MatchBehavior
	SkipPaths              []dt.PathSegment
	SkipDirFunc            SkipEntryFunc
	SkipEntryFunc          SkipEntryFunc
	ParseEntryFunc         ParseEntryFunc
	Logger                 *slog.Logger
	Writer                 writer
}

func NewDirPathScanner(dp dt.DirPath, args DirPathScannerArgs) *DirPathScanner {
	return &DirPathScanner{
		DirPath:                dp,
		ContinueOnErr:          args.ContinueOnErr,
		IncludeNonRegularFiles: args.IncludeNonRegularFiles,
		Writer:                 args.Writer,
		SkipPaths:              args.SkipPaths,
		SkipDirFunc:            args.SkipDirFunc,
		SkipEntryFunc:          args.SkipEntryFunc,
		ParseEntryFunc:         args.ParseEntryFunc,
		MatchBehavior:          args.MatchBehavior,
	}
}

var ErrFileOperation = errors.New("file operation error")

func (ds *DirPathScanner) Scan() (entries []dt.EntryPath, err error) {
	var exists bool
	var dir dt.DirPath

	// Expand tilde, make absolute, and convert string to dt.DirPath

	// TOOD Update this to process all ds.DirPaths elments

	dir, err = ds.DirPath.Normalize()
	if err != nil {
		err = NewErr(ErrFileOperation, err)
		goto end
	}

	// Verify directory exists
	exists, err = dir.Exists()
	if err != nil {
		err = NewErr(dt.ErrNotDirectory, err)
		goto end
	}

	if !exists {
		err = NewErr(dt.ErrFileNotExists)
		goto end
	}

	entries, err = ds.scanDir(dir)
	if err != nil {
		goto end
	}

end:
	if err != nil {
		err = WithErr(err,
			"file_op", "scan_dir",
			"dir", dir,
		)
	}
	return entries, err
}

// scanDir recursively scans a directory for go.mod files
func (ds *DirPathScanner) scanDir(root dt.DirPath) (entries []dt.EntryPath, err error) {
	var de dt.DirEntry

	skipDirs := ds.mapSkipDirs()
	for de, err = range root.Walk() {
		if err != nil {
			// Permission denied and similar errors are non-fatal by default
			// Write to stderr unless --silence-errors is set
			if ds.Writer != nil {
				ds.Writer.Errorf("Error accessing path %s; %v", de.DirPath(), err)
			}
			if ds.Logger != nil {
				ds.Logger.Warn("Error accessing path", "dir_path", de.DirPath(), "error", err)
			}
			// Continue scanning other paths
			continue
		}

		// Skip if not a go.mod file
		if de.IsDir() {
			_, skip := skipDirs[de.Rel.Base()]
			if skip {
				// No need to traverse down some subdirectories
				de.SkipDir()
			}
			if ds.SkipDirFunc != nil && ds.SkipDirFunc(root, &de) {
				continue
			}
		}
		if !ds.IncludeNonRegularFiles && !de.Entry.Type().IsRegular() {
			continue
		}

		if ds.SkipEntryFunc != nil && ds.SkipEntryFunc(root, &de) {
			continue
		}
		{
			var collect bool
			entry := de.EntryPath()
			if ds.ParseEntryFunc != nil {
				entry = ds.ParseEntryFunc(entry, root, de)
			}
			switch ds.MatchBehavior {
			case WriteAndCollectOnMatch:
				collect = true
				fallthrough
			case WriteOnMatch:
				if ds.Writer == nil {
					panic("dtx.DirPathScanner.Scan() called with WriteOnMatch behavior requested but no Writer provided.")
				}
				// Output the go.mod path
				ds.Writer.Printf("%s\n", entry)
			case CollectOnMatch:
				fallthrough
			default:
				collect = true
			}
			if collect {
				entries = append(entries, entry)
			}
		}
	}
	return entries, err
}

// scanDir recursively scans a directory for go.mod files
func (ds *DirPathScanner) mapSkipDirs() (m map[dt.PathSegment]struct{}) {
	m = make(map[dt.PathSegment]struct{}, len(ds.SkipPaths))
	for _, path := range ds.SkipPaths {
		m[path] = struct{}{}
	}
	return m
}

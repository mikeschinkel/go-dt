package dtglob

import (
	"fmt"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/mikeschinkel/go-dt"
)

// CopyTo applies all rules to copy files to the installation directory
func (grs *GlobRules) CopyTo(installDir dt.DirPath, opts *dt.CopyOptions) (err error) {
	var errs []error

	// Normalize opts
	if opts == nil {
		opts = new(dt.CopyOptions)
	}

	// 3. Copy all files (no MkdirAll per file)
	for _, rule := range grs.Rules {
		err = rule.copyTo(grs.BaseDir, installDir, opts)
		if err != nil && !rule.Optional {
			errs = dt.AppendErr(errs, err)
		}
	}

	err = dt.CombineErrs(errs)
	//end:
	return err
}

// collectDestDirs gathers all unique destination directories from all rules
func (grs *GlobRules) collectDestDirs(installDir dt.DirPath) (dirs map[dt.DirPath]struct{}, err error) {
	var rule GlobRule
	var matches []string
	var match string
	var destPath dt.Filepath
	var destDir dt.DirPath
	var errs []error

	dirs = make(map[dt.DirPath]struct{})

	for _, rule = range grs.Rules {
		// Match files for this rule
		matches, _ = doublestar.Glob(grs.BaseDir.DirFS(), string(rule.From))

		for _, match = range matches {
			var status dt.EntryStatus
			destPath, err = rule.computeDestPath(dt.Filepath(match), installDir)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			status, err = destPath.Status()
			if err != nil {
				errs = append(errs, err)
				continue
			}
			switch status {
			case dt.IsFileEntry:
				destDir = destPath.Dir()
			case dt.IsDirEntry, dt.IsMissingEntry:
				destDir = dt.DirPath(destPath)
			default:
				err = dt.NewErr(
					dt.ErrNotFileOrDirectory,
					"entry_status", status,
				)
			}
			dirs[destDir] = struct{}{}
		}
	}
	err = dt.CombineErrs(errs)
	return dirs, err
}

// copyTo processes a single rule
func (rule *GlobRule) copyTo(baseDir, installDir dt.DirPath, opts *dt.CopyOptions) (err error) {
	var matches []string
	var match string
	var sourcePath dt.Filepath
	var destPath dt.Filepath
	var errs []error

	// Find all files matching the glob pattern
	matches, err = doublestar.Glob(baseDir.DirFS(), string(rule.From))
	if err != nil {
		err = fmt.Errorf("glob pattern error for '%s': %w", rule.From, err)
		goto end
	}

	if len(matches) == 0 {
		if rule.Optional {
			goto end // No matches for optional rule is OK
		}
		err = fmt.Errorf("no files matched pattern: %s", rule.From)
		goto end
	}

	// Process each matched file
	for _, match = range matches {
		sourcePath = dt.FilepathJoin(baseDir, match)

		var status dt.EntryStatus
		status, err = sourcePath.Status()
		if err != nil {
			errs = dt.AppendErr(errs, err)
			continue
		}
		if status == dt.IsDirEntry {
			// Skip directories - we only copy files
			continue
		}

		// Determine destination path
		destPath, err = rule.computeDestPath(dt.Filepath(match), installDir)
		if err != nil {
			errs = dt.AppendErr(errs, err)
			continue
		}
		var exists bool
		exists, err = destPath.Dir().Exists()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !exists {
			err = destPath.Dir().MkdirAll(0755)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
		err = sourcePath.CopyTo(destPath, opts)
		if err != nil {
			err = dt.WithErr(err,
				dt.ErrFailedToCopyFile,
				"source", match,
				"destination", destPath,
			)
			errs = dt.AppendErr(errs, err)
			continue
		}
	}

	err = dt.CombineErrs(errs)

end:
	return err
}

// computeDestPath determines the destination path for a matched file
func (rule *GlobRule) computeDestPath(matchedPath dt.Filepath, destDir dt.DirPath) (destPath dt.Filepath, err error) {
	var parts []string
	var basePattern dt.DirPath
	var relativePath dt.RelFilepath

	// If pattern contains **, we need to preserve directory structure
	if rule.From.Contains("**") {
		// Extract the base pattern (part before **)
		parts = rule.From.Split("**")
		basePattern = dt.DirPath(strings.TrimSuffix(parts[0], "/"))

		// Remove the base pattern from the matched path to get the relative part
		relativePath, err = matchedPath.Rel(basePattern)
		if err != nil {
			goto end
		}

		// If To is "." or ends with /, append the relative path
		if rule.To == "." || rule.To.HasSuffix("/") {
			destPath = dt.FilepathJoin3(destDir, rule.To, relativePath)
			goto end
		}

		// Otherwise, use To as exact destination
		destPath = dt.FilepathJoin(destDir, rule.To)
		goto end
	}

	// For simple patterns without **, check if To is a directory or file
	if rule.To.HasSuffix("/") {
		// To is a directory - use the filename from the match
		destPath = dt.FilepathJoin3(destDir, rule.To, matchedPath.Base())
		goto end
	}

	// To is a specific file path
	destPath = dt.FilepathJoin(destDir, rule.To)

end:
	return destPath, err
}

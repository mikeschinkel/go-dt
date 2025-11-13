package dtglob

import (
	"fmt"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/mikeschinkel/go-dt"
)

// CopyTo applies all rules to copy files to the installation directory
func (grs *GlobRules) CopyTo(installDir dt.DirPath, opts *dt.CopyOptions) (err error) {
	var errs []error
	var dirs map[dt.DirPath]bool
	var dir dt.DirPath

	// Normalize opts
	if opts == nil {
		opts = new(dt.CopyOptions)
	}

	// 1. Collect all unique destination directories
	dirs = grs.collectDestDirs(installDir)

	// 2. Create all directories once
	for dir = range dirs {
		err = dir.MkdirAll(0755)
		errs = dt.AppendErr(errs, err)
	}

	// 3. Copy all files (no MkdirAll per file)
	for _, rule := range grs.Rules {
		err = rule.copyTo(grs.BaseDir, installDir, opts)
		if err != nil && !rule.Optional {
			errs = dt.AppendErr(errs, err)
		}
	}

	err = dt.CombineErrs(errs)
	return err
}

// collectDestDirs gathers all unique destination directories from all rules
func (grs *GlobRules) collectDestDirs(installDir dt.DirPath) (dirs map[dt.DirPath]bool) {
	var rule GlobRule
	var matches []string
	var match string
	var destPath dt.Filepath
	var destDir dt.DirPath

	dirs = make(map[dt.DirPath]bool)

	for _, rule = range grs.Rules {
		// Match files for this rule
		matches, _ = doublestar.Glob(grs.BaseDir.DirFS(), string(rule.From))

		for _, match = range matches {
			destPath, _ = rule.computeDestPath(dt.Filepath(match), installDir)
			destDir = destPath.Dir()
			dirs[destDir] = true
		}
	}

	return dirs
}

// copyTo processes a single rule
func (rule *GlobRule) copyTo(baseDir, installDir dt.DirPath, opts *dt.CopyOptions) (err error) {
	var matches []string
	var match string
	var sourcePath dt.Filepath
	var info os.FileInfo
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

		// Skip directories - we only copy files
		info, err = sourcePath.Stat()
		if err != nil {
			errs = dt.AppendErr(errs, err)
			continue
		}
		if info.IsDir() {
			continue
		}

		// Determine destination path
		destPath, err = rule.computeDestPath(dt.Filepath(match), installDir)
		if err != nil {
			errs = dt.AppendErr(errs, err)
			continue
		}

		// Perform the copy
		err = sourcePath.CopyTo(destPath, opts)
		if err != nil {
			err = fmt.Errorf("failed to copy %s to %s: %w", match, destPath, err)
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

		// If To ends with /, append the relative path
		if rule.To.HasSuffix("/") {
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

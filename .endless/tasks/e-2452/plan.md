In .github/workflows/release.yml, step 'Generate release notes', add ':(exclude).endless' to root's PATHS array, after the nested-module exclusions. git's pathspec then drops any commit whose changes are all under .endless/ and keeps any commit that also changes something else. Nested modules' PATHS (-- "$DIR") is unchanged; .endless/ never falls inside a nested module's directory.

Add a verify suite at .endless/tasks/e-2452/verify.sh, sourcing _harness.sh, that greps release.yml for the new exclusion and runs root's git log with release.yml's PATHS over v0.6.0..v0.6.1, checking that no listed commit touches only .endless/ and that be26016's change (PathSegments.IsRoot) is listed.

After landing, with Mike's confirmation, regenerate v0.6.1's GitHub release body with the new PATHS and replace it with `gh release edit v0.6.1 --notes-file`, keeping its heading, description and Full Changelog line.

Other go-pkgs repos get the same fix through Gomion's E-2454, not here.

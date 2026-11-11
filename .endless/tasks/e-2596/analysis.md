Every explicit value equals the value iota assigns today, so callers, String(), GetEntryStatus and any stored value see no change; this is not a compatibility break. The one real risk is a transcription mistake, which a test pinning each constant to its number catches.

DontFollowSymlinks becomes `1 << 0` rather than `1`, so the declaration still reads as a bit flag and the next flag is obviously `1 << 1`.

go-dt's existing explicit-value style is the model: TildeOption in tilde.go and dtx's MatchBehavior give each constant a literal value.

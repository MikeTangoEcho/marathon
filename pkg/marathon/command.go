package marathon

import "regexp"

const CommandPrefix string = "!"
const PlayCommandPattern string = "^" + CommandPrefix + "play (?P<path>\\S+)$"
const ShutdownCommandPattern string = "^" + CommandPrefix + "shutdown$"

var PlayCommandRegexp *regexp.Regexp = regexp.MustCompile(PlayCommandPattern)
var ShutdownCommandRegexp *regexp.Regexp = regexp.MustCompile(ShutdownCommandPattern)

package util

import (
    "regexp"
    "strings"
)

// GroupsFromRegex gets RegExp groups after matching it on a line
func GroupsFromRegex(rx, line string) map[string]string {
    // We remove tabs and newlines from the regex. This is purely cosmetic,
    // as the regex input can be quite long and it's nice for the caller to
    // be able to format it in a more readable way.
    rx = strings.Replace(rx, "\r", "", -1)
    rx = strings.Replace(rx, "\n", "", -1)
    rx = strings.Replace(rx, "\t", "", -1)
    re := regexp.MustCompile(rx)

    match := re.FindStringSubmatch(line)
    if len(match) == 0 {
        return nil
    }

    result := make(map[string]string)
    for i, name := range re.SubexpNames() {
        if i != 0 {
            result[name] = match[i]
        }
    }

    return result
}

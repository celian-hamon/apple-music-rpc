tell application "Music"
    set rawData to (get raw data of artwork 1 of current track)
end tell
    set newPath to ("Macintosh HD:Users:cece:Code:go-apple-music-rpc:tmp.jpg") as text
    tell me to set fileRef to (open for access newPath with write permission)
    write rawData to fileRef starting at 0
    tell me to close access fileRef
    return "DONE"

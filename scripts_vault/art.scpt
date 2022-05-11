tell application "Music"
    set rawData to (get raw data of artwork 1 of current track)
end tell
tell application "Finder"
    set current_path to container of (((path to me as text) & "::") as alias) as string
end tell
set newPath to ((current_path as text) & "tmp.jpg") as text
tell me to set fileRef to (open for access newPath with write permission)
write rawData to fileRef starting at 0
tell me to close access fileRef
return "DONE"

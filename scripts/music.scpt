if application "Music" is running then
    tell application "Music"
        if player state is playing or player state is paused then
            set currentTrack to current track

            return {get player state} & {get artist of currentTrack} & {get name of currentTrack} & {get album of currentTrack} & {get kind of currentTrack} & {get duration of currentTrack} & {player position} & {get genre of current track} & {get id of current track} 
        else
            return "stopped"
        end if
    end tell
else
    return "stopped"
end if
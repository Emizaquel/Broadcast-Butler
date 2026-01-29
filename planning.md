Rules (config)

automations

Schedule

# OBS hookup
- switch scenes
- update element
- 


# mod interface
- user tracker
 - ban/timeout
 - Suspicious (local for yt/linked to twitch API)
 - bond accounts (connecting known cross platform identities)

# Web automations
- I want to be able to automate the report process (at least for bots ig)

# Chat
- message
    - message parts (text/image)
- delete/timeout/ban/

# Stream Events
- Twitch: Hype Train
- Twitch bits /YT superchat
- YT Sub/Twitch follow
- YT Member/Twitch Sub
- YT/Twitch: Charity doation goal?

(assume lower information in the compiled thingamajig)
Platform
- channel
    - name
    - id
    - auth token
    - update event handling
        - stream name
        - stream category
        - schedule
- users (table in db thing?)
    - profile image url
    - type (normal/subscriber/paid/mod/owner)
    - display name
    - true name
    - verified
    - messages (table)
        - timestamp
        - superchat / cheers ig
        - parts
            - text or image
            - content (content of text/url for image)
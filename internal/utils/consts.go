package utils

const (
	CommandAlbumsTitle = "Download albums"
	CommandAlbumsDesc  = "Download photos from albums"

	CommandGroupsTitle = "Download groups"
	CommandGroupsDesc  = "Download posts from groups"

	CommandAuthTitle = "Authenticate"
	CommandAuthDesc  = "Authenticate with VK"

	CommandTokenTitle = "Save token"
	CommandTokenDesc  = "Save app token"

	CommandHelpTitle = "Help"
	CommandHelpDesc  = "How to use VKscape"

	MenuQuit = "Quit"

	AppHelpText = `VKscape is a CLI/TUI tool for downloading your VK archive.

Authentication:
- App token (--token): Can download public content. Private albums are not available.
- User auth (--user): Can access your own private content.

Notes:
- For albums: use --owner and optional --ids.
- For groups: use --ids with negative IDs or handles.`
)

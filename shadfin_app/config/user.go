package config

import "time"

type UserStore struct {
	User *User `json:"user"`
}

type User struct {
	Name                      string     `json:"Name"`
	ServerID                  string     `json:"ServerId"`
	ID                        string     `json:"Id"`
	HasPassword               bool       `json:"HasPassword"`
	HasConfiguredPassword     bool       `json:"HasConfiguredPassword"`
	HasConfiguredEasyPassword bool       `json:"HasConfiguredEasyPassword"`
	EnableAutoLogin           bool       `json:"EnableAutoLogin"`
	LastLoginDate             *time.Time `json:"LastLoginDate"`
	LastActivityDate          *time.Time `json:"LastActivityDate"`
	Configuration             struct {
		AudioLanguagePreference    string   `json:"AudioLanguagePreference"`
		PlayDefaultAudioTrack      bool     `json:"PlayDefaultAudioTrack"`
		SubtitleLanguagePreference string   `json:"SubtitleLanguagePreference"`
		DisplayMissingEpisodes     bool     `json:"DisplayMissingEpisodes"`
		GroupedFolders             []string `json:"GroupedFolders"`
		SubtitleMode               string   `json:"SubtitleMode"`
		DisplayCollectionsView     bool     `json:"DisplayCollectionsView"`
		EnableLocalPassword        bool     `json:"EnableLocalPassword"`
		OrderedViews               []string `json:"OrderedViews"`
		LatestItemsExcludes        []string `json:"LatestItemsExcludes"`
		MyMediaExcludes            []string `json:"MyMediaExcludes"`
		HidePlayedInLatest         bool     `json:"HidePlayedInLatest"`
		RememberAudioSelections    bool     `json:"RememberAudioSelections"`
		RememberSubtitleSelections bool     `json:"RememberSubtitleSelections"`
		EnableNextEpisodeAutoPlay  bool     `json:"EnableNextEpisodeAutoPlay"`
		CastReceiverID             string   `json:"CastReceiverId"`
	} `json:"Configuration"`
	Policy struct {
		IsAdministrator                  bool     `json:"IsAdministrator"`
		IsHidden                         bool     `json:"IsHidden"`
		EnableCollectionManagement       bool     `json:"EnableCollectionManagement"`
		EnableSubtitleManagement         bool     `json:"EnableSubtitleManagement"`
		EnableLyricManagement            bool     `json:"EnableLyricManagement"`
		IsDisabled                       bool     `json:"IsDisabled"`
		BlockedTags                      []string `json:"BlockedTags"`
		AllowedTags                      []string `json:"AllowedTags"`
		EnableUserPreferenceAccess       bool     `json:"EnableUserPreferenceAccess"`
		AccessSchedules                  []string `json:"AccessSchedules"`
		BlockUnratedItems                []string `json:"BlockUnratedItems"`
		EnableRemoteControlOfOtherUsers  bool     `json:"EnableRemoteControlOfOtherUsers"`
		EnableSharedDeviceControl        bool     `json:"EnableSharedDeviceControl"`
		EnableRemoteAccess               bool     `json:"EnableRemoteAccess"`
		EnableLiveTvManagement           bool     `json:"EnableLiveTvManagement"`
		EnableLiveTvAccess               bool     `json:"EnableLiveTvAccess"`
		EnableMediaPlayback              bool     `json:"EnableMediaPlayback"`
		EnableAudioPlaybackTranscoding   bool     `json:"EnableAudioPlaybackTranscoding"`
		EnableVideoPlaybackTranscoding   bool     `json:"EnableVideoPlaybackTranscoding"`
		EnablePlaybackRemuxing           bool     `json:"EnablePlaybackRemuxing"`
		ForceRemoteSourceTranscoding     bool     `json:"ForceRemoteSourceTranscoding"`
		EnableContentDeletion            bool     `json:"EnableContentDeletion"`
		EnableContentDeletionFromFolders []string `json:"EnableContentDeletionFromFolders"`
		EnableContentDownloading         bool     `json:"EnableContentDownloading"`
		EnableSyncTranscoding            bool     `json:"EnableSyncTranscoding"`
		EnableMediaConversion            bool     `json:"EnableMediaConversion"`
		EnabledDevices                   []string `json:"EnabledDevices"`
		EnableAllDevices                 bool     `json:"EnableAllDevices"`
		EnabledChannels                  []string `json:"EnabledChannels"`
		EnableAllChannels                bool     `json:"EnableAllChannels"`
		EnabledFolders                   []string `json:"EnabledFolders"`
		EnableAllFolders                 bool     `json:"EnableAllFolders"`
		InvalidLoginAttemptCount         int      `json:"InvalidLoginAttemptCount"`
		LoginAttemptsBeforeLockout       int      `json:"LoginAttemptsBeforeLockout"`
		MaxActiveSessions                int      `json:"MaxActiveSessions"`
		EnablePublicSharing              bool     `json:"EnablePublicSharing"`
		BlockedMediaFolders              []string `json:"BlockedMediaFolders"`
		BlockedChannels                  []string `json:"BlockedChannels"`
		RemoteClientBitrateLimit         int      `json:"RemoteClientBitrateLimit"`
		AuthenticationProviderID         string   `json:"AuthenticationProviderId"`
		PasswordResetProviderID          string   `json:"PasswordResetProviderId"`
		SyncPlayAccess                   string   `json:"SyncPlayAccess"`
	} `json:"Policy"`
}

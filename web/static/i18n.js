const LANGUAGE_STORAGE_KEY = 'uc_framework_language';
const DEFAULT_LANGUAGE = 'en';

const TRANSLATIONS = {
    en: {
        'title.dashboard': 'UC-Framework Dashboard',
        'title.botControl': 'Bot Control',
        'title.counter': 'Counter Settings',
        'title.plugins': 'Plugins',
        'title.ts3Settings': 'TeamSpeak Settings',
        'title.ts3Connection': 'TS3 Connection Test',
        'title.framework': 'Framework',
        'title.discord': 'Discord Settings',
        'title.announcement': 'Announcements',
        'title.support': 'Support Channels',
        'title.authSettings': 'Auth Settings',
        'title.afkMover': 'AFK Mover',
        'nav.dashboard': 'Dashboard',
        'nav.teamspeak': 'TeamSpeak',
        'nav.discord': 'Discord',
        'nav.ts3Connection': 'TS3 Connection',
        'nav.bot': 'Bot',
        'nav.framework': 'Framework',
        'nav.counter': 'Counter',
        'nav.announcement': 'Announcements',
        'nav.support': 'Support',
        'nav.afkMover': 'AFK Mover',
        'nav.auth': 'Auth',
        'nav.plugins': 'Plugins',
        'nav.settings': 'Settings',
        'nav.logout': 'Logout',
        'nav.language': 'Language',
        'nav.languageEnglish': 'English',
        'nav.languageGerman': 'German',
        'dashboard.heading': 'UC-Dashboard',
        'dashboard.stats': 'Statistics',
        'dashboard.plugins': 'Plugin Status',
        'dashboard.settings': 'Settings',
        'dashboard.admins': 'Admins',
        'dashboard.users': 'Users',
        'dashboard.supportersOnline': 'Supporters online',
        'dashboard.membersOnline': 'Members online',
        'dashboard.loading': 'Loading...',
        'dashboard.statsLoadError': 'Failed to load statistics.',
        'dashboard.settingsLoadError': 'Failed to load settings.',
        'dashboard.settingsUnavailable': 'No settings available.',
        'dashboard.ts3Host': 'TS3 host',
        'dashboard.ts3QueryPort': 'TS3 query port',
        'dashboard.ts3Nickname': 'Bot nickname',
        'dashboard.ts3DefaultChannel': 'Default channel',
        'dashboard.botStatus': 'Bot Status',
        'dashboard.botRunning': 'Running',
        'dashboard.botStopped': 'Stopped',
        'dashboard.botStatusLoadError': 'Failed to load bot status.',
        'plugins.heading': 'Plugins',
        'plugins.description': 'Enable or disable plugins here. Detailed settings remain on each plugin page.',
        'plugins.listHeading': 'Plugin Management',
        'plugin.adminCounter.name': 'Admin Counter',
        'plugin.adminCounter.description': 'Counts online admins based on configured server groups.',
        'plugin.memberCounter.name': 'Member Counter',
        'plugin.memberCounter.description': 'Counts online members and excludes configured groups or nicknames.',
        'plugin.afkMover.name': 'AFK Mover',
        'plugin.afkMover.description': 'Moves inactive users to the AFK channel after the configured timeout.',
        'plugin.combinedStats.name': 'Combined Stats',
        'plugin.combinedStats.description': 'Provides combined admin and member statistics for the dashboard and API.',
        'plugin.supportControl.name': 'Support Channel Control',
        'plugin.supportControl.description': 'Opens and closes support channels via permissions, with optional server messages and join pokes plus auto-scheduling.',
        'common.yes': 'Yes',
        'common.no': 'No',
        'common.enabled': 'Enabled',
        'common.provider': 'Provider',
        'common.user': 'User',
        'common.active': 'active',
        'common.inactive': 'inactive',
        'common.saved': 'Saved!',
        'common.error': 'Error',
        'common.loading': 'Loading...',
        'common.close': 'Close',
        'common.remove': 'Remove',
        'common.noEntries': 'No entries.',
        'common.channelSelectionConnectedHint': 'Channel selection is only possible when the bot is connected to TeamSpeak.',
        'framework.heading': 'Framework',
        'framework.description': 'Here you can see framework metadata, choose the active communication platform setup, and restart the complete framework process.',
        'framework.platformTitle': 'Platform Setup',
        'framework.platformHint': 'Choose whether the framework should be prepared for TeamSpeak, Discord, or both platforms together.',
        'framework.platformMode': 'Communication platforms',
        'framework.platformTeamspeak': 'TeamSpeak',
        'framework.platformDiscord': 'Discord',
        'framework.platformBoth': 'TeamSpeak and Discord',
        'framework.savePlatforms': 'Save platform setup',
        'framework.settingsLoaded': 'Framework platform settings loaded.',
        'framework.settingsLoadError': 'Failed to load framework settings: {error}',
        'framework.saved': 'Framework platform settings saved.',
        'framework.saveError': 'Failed to save framework settings: {error}',
        'framework.infoTitle': 'Framework Information',
        'framework.name': 'Name',
        'framework.version': 'Version',
        'framework.updatedAt': 'Updated at',
        'framework.latestVersion': 'Latest version',
        'framework.upToDate': 'Up to date',
        'framework.infoLoaded': 'Framework information loaded.',
        'framework.infoLoadError': 'Failed to load framework information: {error}',
        'framework.restartTitle': 'Framework Restart',
        'framework.restartHint': 'Restarts the entire framework process, including web API and runtime services. (Does not work on all systems.)',
        'framework.restartButton': 'Restart framework',
        'framework.restartConfirm': 'Do you really want to restart the complete framework now?',
        'framework.restartScheduled': 'Framework restart has been scheduled. This page will reload shortly.',
        'framework.restartFailed': 'Framework restart failed.',
        'framework.restartError': 'Failed to restart framework: {error}',
        'announcement.heading': 'Server Announcements',
        'announcement.description': 'Send announcements to all connected users. Configure the announcement message and set it to send once, at a fixed time, or at regular intervals.',
        'announcement.settingsTitle': 'Announcement Message',
        'announcement.message': 'Message',
        'announcement.messagePlaceholder': 'Type your announcement here...',
        'announcement.messageEmpty': 'Please enter a message first.',
        'announcement.messageTooLong': 'The message must not be longer than 500 characters.',
        'announcement.scheduleTitle': 'Scheduling',
        'announcement.repeatTitle': 'Repeat announcement',
        'announcement.repeatHint': 'Send the announcement automatically at a fixed or recurring schedule. Disabled = send once only.',
        'announcement.scheduleMode': 'Schedule mode',
        'announcement.modeOnce': 'Send once now',
        'announcement.modeInterval': 'Repeat every X minutes',
        'announcement.modeTime': 'Send at a specific time',
        'announcement.intervalTime': 'Interval (minutes)',
        'announcement.intervalHint': 'Minimum 10 minutes',
        'announcement.intervalCount': 'How many times to send',
        'announcement.fixedTime': 'Daily send time',
        'announcement.invalidTime': 'Invalid time format. Use HH:MM.',
        'announcement.save': 'Save announcement settings',
        'announcement.saved': 'Announcement settings saved.',
        'announcement.loaded': 'Announcement settings loaded.',
        'announcement.manualTitle': 'Send now',
        'announcement.manualHint': 'Send the current announcement message immediately to all connected users.',
        'announcement.lastSent': 'Last sent',
        'announcement.neverSent': 'Never sent yet',
        'announcement.sendNow': 'Send announcement now',
        'announcement.sent': 'Announcement sent successfully.',
        'announcement.refresh': 'Refresh status',
        'announcement.settingsLoadError': 'Failed to load announcement settings: {error}',
        'announcement.statusLoadError': 'Failed to load announcement status: {error}',
        'announcement.saveError': 'Failed to save announcement settings: {error}',
        'announcement.sendError': 'Failed to send announcement: {error}',
        'bot.heading': 'Bot Control',
        'bot.description': 'Here you can start, stop, and restart the bot. The watchdog automatically restarts the bot if it crashes or a start attempt fails, and retries every 60-120 seconds.',
        'bot.status': 'Bot Status',
        'bot.running': 'Running',
        'bot.pid': 'PID',
        'bot.desiredState': 'Desired state',
        'bot.lastAction': 'Last action',
        'bot.lastError': 'Last error',
        'bot.start': 'Start bot',
        'bot.stop': 'Stop bot',
        'bot.restart': 'Restart bot',
        'bot.watchdog': 'Watchdog',
        'bot.watchdogRunning': 'Watchdog running',
        'bot.watchdogPid': 'Watchdog PID',
        'bot.startWatchdog': 'Start watchdog',
        'bot.stopWatchdog': 'Stop watchdog',
        'bot.refreshStatus': 'Refresh status',
        'bot.system': 'System',
        'bot.os': 'Operating system',
        'bot.arch': 'Architecture',
        'bot.workingDir': 'Working directory',
        'bot.botFile': 'Bot executable',
        'bot.supervisorScript': 'Supervisor script',
        'bot.watchdogScript': 'Watchdog script',
        'bot.logFile': 'Log file',
        'bot.logs': 'Log Viewer',
        'bot.lastLines': 'Last lines',
        'bot.filter': 'Filter',
        'bot.loadLogs': 'Load logs',
        'bot.noLogLines': 'No log lines available.',
        'bot.logsLoaded': 'Log loaded ({lines} lines, filter: {level}).',
        'bot.actionExecuted': 'Action {action} executed.',
        'bot.watchdogActionExecuted': 'Watchdog {action} executed.',
        'bot.statusLoadError': 'Failed to load status: {error}',
        'bot.logsLoadError': 'Failed to load logs: {error}',
        'log.level.all': 'All',
        'log.level.error': 'Error',
        'log.level.warn': 'Warn',
        'log.level.info': 'Info',
        'log.level.debug': 'Debug',
        'counter.heading': 'Counter Settings',
        'counter.adminCounter': 'Admin Counter',
        'counter.adminDescription': 'Server group IDs that count as admins.',
        'counter.memberCounter': 'Member Counter',
        'counter.memberDescription': 'All users are counted except for the following exclusions.',
        'counter.excludedGroups': 'Excluded groups',
        'counter.excludedGroupsDesc': 'Users with these server group IDs are not counted.',
        'counter.excludedNicknames': 'Excluded nicknames',
        'counter.excludedNicknamesDesc': 'Users with exactly these nicknames are not counted.',
        'counter.groupId': 'Group ID',
        'counter.nickname': 'Nickname',
        'counter.renameSection': 'Channel rename on count change',
        'counter.renameHint': 'Use a custom token. Example: template "cspacer7 %%" and token "%%".',
        'counter.renameChannelId': 'Channel ID to rename',
        'counter.chooseRenameChannel': '⌬',
        'counter.renameTemplate': 'New channel name template',
        'counter.renameToken': 'Count placeholder token (e.g. {count} or %%)',
        'counter.channelPopupTitle': 'Available channels',
        'counter.channelPopupHint': 'Select a channel. The selected channel ID is saved immediately.',
        'counter.add': 'Add',
        'counter.saveAdmin': 'Save Admin Counter',
        'counter.saveMember': 'Save Member Counter',
        'counter.pluginNotLoaded': 'Plugin is inactive. Activation happens on the plugins page.',
        'counter.saveVerifyFailed': 'Saved response does not contain your rename fields. Reload the page and retry.',
        'ts3.heading': 'TeamSpeak Settings',
        'ts3.description': 'Here you can edit all TS3 connection and bot data. After saving, a bot restart is required for the new connection to become active.',
        'ts3.server': 'Server',
        'ts3.host': 'TS3 host address',
        'ts3.queryPort': 'TS3 query port',
        'ts3.voicePort': 'TS3 voice port',
        'ts3.queryLogin': 'Query Login',
        'ts3.queryUser': 'TS3 query user',
        'ts3.queryPassword': 'TS3 query password',
        'ts3.botBehavior': 'Bot Behavior',
        'ts3.botNickname': 'Bot nickname',
        'ts3.defaultChannel': 'Default channel',
        'ts3.querySlowmode': 'Query slowmode (ms)',
        'ts3.chooseChannel': '⌬',
        'ts3.channelPopupTitle': 'Available channels',
        'ts3.channelPopupHint': 'Select a channel to use as default channel.',
        'ts3.channelLoadError': 'Channels could not be loaded.',
        'ts3.noChannels': 'No channels available.',
        'ts3.unnamedChannel': 'Unnamed channel',
        'ts3.save': 'Save TS3 settings',
        'ts3.loadError': 'Failed to load TS3 settings.',
        'ts3.loadException': 'Error while loading TS3 settings.',
        'ts3.savedRestart': 'Saved. Please restart the bot so the changes take effect.',
        'ts3.saved': 'Saved.',
        'discord.heading': 'Discord Settings',
        'discord.description': 'Prepare Discord integration completely in the web interface. All important bot, server, channel, and role settings can be maintained here.',
        'discord.topHelpHtml': '<p><strong>Discord setup quick help</strong></p><p>Bot token: secret key that allows this framework to log in as your Discord bot.</p><p>Application ID: unique ID of your Discord application (Developer Portal).</p><p>Guild ID: unique ID of the Discord server where the bot should run.</p>',
        'discord.enableTitle': 'Discord integration',
        'discord.enableHint': 'Enable this as soon as the Discord connection should be configured and prepared in the framework.',
        'discord.authTitle': 'Bot Access',
        'discord.authHelpHtml': '<strong>Bot Access</strong><p>Bot token: secret key for bot login.</p><p>Application ID: unique ID of the Discord application.</p><p>Guild ID: unique ID of your Discord server.</p>',
        'discord.botToken': 'Bot token',
        'discord.botTokenPlaceholder': 'Discord bot token',
        'discord.applicationId': 'Application ID',
        'discord.applicationIdPlaceholder': 'Discord application ID',
        'discord.guildId': 'Guild ID',
        'discord.guildIdPlaceholder': 'Discord guild ID',
        'discord.behaviorTitle': 'Bot Behavior',
        'discord.afkTitle': 'Discord AFK Handling',
        'discord.afkKickTitle': 'Kick inactive voice users',
        'discord.afkKickHint': 'Discord-specific: users who have not spoken or changed their voice state for the configured duration are disconnected from voice.',
        'discord.afkInactivityMinutes': 'Inactivity timeout (minutes)',
        'discord.botDisplayName': 'Display name',
        'discord.botDisplayNamePlaceholder': 'UC-Framework',
        'discord.statusText': 'Status text',
        'discord.statusTextPlaceholder': 'Watching support activity',
        'discord.commandPrefix': 'Command prefix',
        'discord.commandPrefixPlaceholder': '!',
        'discord.channelsTitle': 'Channels',
        'discord.channelIdPlaceholder': 'Discord channel ID',
        'discord.chooseChannel': 'Choose channel',
        'discord.logChannelId': 'Log channel ID',
        'discord.announcementChannelId': 'Announcement channel ID',
        'discord.supportCategoryId': 'Support category ID',
        'discord.supportLogChannelId': 'Support log channel ID',
        'discord.channelPopupTitle': 'Select Discord channel',
        'discord.channelPopupHint': 'Click one channel to use its ID for the selected field.',
        'discord.channelsLoadError': 'Discord channels could not be loaded.',
        'discord.noChannels': 'No Discord channels available.',
        'discord.unnamedChannel': 'Unnamed channel',
        'discord.rolesTitle': 'Roles',
        'discord.adminRoles': 'Admin role IDs',
        'discord.supporterRoles': 'Supporter role IDs',
        'discord.botRoles': 'Bot role IDs',
        'discord.chooseRole': 'Choose roles',
        'discord.roleIdPlaceholder': 'Discord role ID',
        'discord.rolePopupTitle': 'Select Discord roles',
        'discord.rolePopupHint': 'Click roles to add or remove them. Multiple selection is possible.',
        'discord.rolesLoadError': 'Discord roles could not be loaded.',
        'discord.noRoles': 'No Discord roles available.',
        'discord.unnamedRole': 'Unnamed role',
        'discord.addRole': 'Add',
        'discord.noRolesSelected': 'No roles selected.',
        'discord.save': 'Save Discord settings',
        'discord.loaded': 'Discord settings loaded.',
        'discord.loadError': 'Failed to load Discord settings: {error}',
        'discord.saved': 'Discord settings saved.',
        'discord.savedRestart': 'Saved. Please restart the bot so the changes take effect.',
        'discord.saveError': 'Failed to save Discord settings: {error}',
        'ts3conn.heading': 'TS3 Connection Test',
        'ts3conn.description': 'Check whether the bot is currently connected to TeamSpeak and run an active connection test.',
        'ts3conn.currentStatus': 'Current status',
        'ts3conn.connected': 'Connected',
        'ts3conn.host': 'Host',
        'ts3conn.port': 'Query Port',
        'ts3conn.lastCheck': 'Last check',
        'ts3conn.lastError': 'Last error',
        'ts3conn.botDisabled': 'Disabled',
        'ts3conn.botDisabledInfo': 'The bot is disabled.',
        'ts3conn.refresh': 'Refresh status',
        'ts3conn.runTest': 'Run connection test',
        'ts3conn.statusLoaded': 'Status loaded.',
        'ts3conn.statusLoadError': 'Status could not be loaded: {error}',
        'ts3conn.testOk': 'Connection test successful.',
        'ts3conn.testFailed': 'Connection test failed.',
        'ts3conn.testError': 'Connection test failed: {error}',
        'auth.heading': 'Authentication',
        'auth.forcePasswordTitle': 'Password change required',
        'auth.forcePasswordText': 'A new password must be set on first login. Other areas stay locked until then.',
        'auth.currentMode': 'Current mode',
        'auth.loginMode': 'Login mode',
        'auth.activeMode': 'Active login mode',
        'auth.modeNone': 'Without login',
        'auth.modeLocal': 'Local login',
        'auth.modeRanksystem': 'TSN-Ranksystem login',
        'auth.modeLocalRanksystem': 'Local + TSN-Ranksystem login',
        'auth.saveMode': 'Save login mode',
        'auth.modeSaved': 'Login mode saved.',
        'auth.noLoginWarningTitle': 'Security warning',
        'auth.noLoginWarningText': 'Login is disabled. This is insecure and should be changed immediately.',
        'auth.rankHealth': 'Rank-System Health',
        'auth.rankHealthText': 'Checks whether the configured rank-system login endpoint is reachable.',
        'auth.startHealth': 'Start health check',
        'auth.changeLocalPassword': 'Change local password',
        'auth.currentPassword': 'Current password',
        'auth.newPassword': 'New password',
        'auth.savePassword': 'Save password',
        'auth.changeLocalUsername': 'Change local username',
        'auth.newUsername': 'New username',
        'auth.saveUsername': 'Save username',
        'auth.modeLoadError': 'Failed to load auth mode.',
        'auth.healthFailed': 'Health check failed.',
        'auth.settingsLoadError': 'Error loading auth settings.',
        'auth.usernameTooShort': 'New username must be at least 3 characters long.',
        'auth.usernameChanged': 'Username changed successfully.',
        'auth.passwordTooShort': 'New password must be at least 8 characters long.',
        'auth.passwordChanged': 'Password changed successfully.',
        'login.subtitle': 'Web interface login',
        'login.text': 'Please sign in with your credentials.',
        'login.username': 'Username',
        'login.password': 'Password',
        'login.signIn': 'Sign in',
        'login.tooManyAttempts': 'Too many failed attempts. Please try again in {seconds} seconds.',
        'login.tooManyAttemptsGeneric': 'Too many failed attempts. Please try again later.',
        'login.failed': 'Login failed. Check your credentials.',
        'login.autoLogout': 'You were logged out automatically after 30 minutes of inactivity.',
        'afk.heading': 'AFK Mover',
        'afk.general': 'General',
        'afk.generalText': 'Plugin activation is managed centrally on the plugins page. This page is for configuration only.',
        'afk.channelId': 'AFK channel ID',
        'afk.channelIdPlaceholder': 'Channel ID',
        'afk.chooseTargetChannel': '⌬',
        'afk.targetPreviewLabel': 'Selected target channel',
        'afk.channelPopupTitle': '⌬',
        'afk.channelPopupHint': 'Select the channel where inactive users should be moved.',
        'afk.chooseExcludedChannel': '⌬',
        'afk.excludedChannelPopupTitle': 'Select excluded channel',
        'afk.excludedChannelPopupHint': 'Select a channel that should be excluded from AFK moving.',
        'afk.channelLoadError': 'Channels could not be loaded.',
        'afk.noChannels': 'No channels available.',
        'afk.unnamedChannel': 'Unnamed channel',
        'afk.timeoutMinutes': 'Inactivity time (minutes)',
        'afk.returnOnActivityTitle': 'Auto return after activity',
        'afk.returnOnActivityLabel': 'Move users back when active again',
        'afk.returnOnActivityHint': 'If enabled, users moved by AFK Mover are returned to their previous channel when they become active again.',
        'afk.excludedChannels': 'Excluded channels',
        'afk.excludedChannelsText': 'Users in these channels are not moved into the AFK channel.',
        'afk.add': 'Add',
        'afk.save': 'Save settings',
        'afk.pluginDisabledInfo': 'The plugin is currently disabled. Activation happens on the plugins page.',
        'afk.noExceptions': 'No exceptions configured.',
        'afk.loadError': 'Failed to load configuration: {error}',
        'support.heading': 'Support Channel Control',
        'support.description': 'Configure support channels, server messages, join pokes, waiting area, and automatic open/close times. Open/close is enforced via channel permissions.',
        'support.settingsTitle': 'Support Settings',
        'support.channels': 'Support channels to open/close',
        'support.waitingArea': 'Support waiting area',
        'support.noWaitingArea': 'No waiting area selected',
        'support.chooseWaitingArea': '⌬',
        'support.waitingPreviewLabel': 'Selected waiting area',
        'support.waitingPopupTitle': 'Select waiting area',
        'support.waitingPopupHint': 'Choose one channel as the waiting area.',
        'support.openPoke': 'Server message when support is opened',
        'support.closePoke': 'Server message when support is closed',
        'support.joinOpenPoke': 'Poke for users joining while support is open',
        'support.joinClosedPoke': 'Poke for users joining while support is closed',
        'support.joinOpenPokePlaceholder': 'Welcome, support is currently open.',
        'support.joinClosedPokePlaceholder': 'Support is currently closed.',
        'support.supporterPoke': 'Poke for supporters when a user joins (open support only)',
        'support.supporterPokePlaceholder': 'New support request from {user}.',
        'support.supporterGroups': 'Supporter server groups (multiple selection)',
        'support.chooseGroups': '👥',
        'support.groupIdPlaceholder': 'Group ID',
        'support.groupPopupTitle': 'Select supporter groups',
        'support.groupPopupHint': 'Click groups to add or remove them. Multiple selection possible.',
        'support.noGroupsSelected': 'No groups selected',
        'support.openPokePlaceholder': 'Support is now open.',
        'support.closePokePlaceholder': 'Support is now closed.',
        'support.autoTitle': 'Automatic opening and closing',
        'support.autoHint': 'Automatically toggle support channels by time (local server time).',
        'support.openTime': 'Open time',
        'support.closeTime': 'Close time',
        'support.save': 'Save support settings',
        'support.saved': 'Support settings saved.',
        'support.loaded': 'Support settings loaded.',
        'support.manualTitle': 'Manual open and close',
        'support.currentState': 'Current state',
        'support.lastAction': 'Last action',
        'support.lastError': 'Last error',
        'support.openNow': 'Open now',
        'support.closeNow': 'Close now',
        'support.refresh': 'Refresh status',
        'support.stateOpen': 'Open',
        'support.stateClosed': 'Closed',
        'support.noChannels': 'No channels available.',
        'support.noChannelsSelected': 'No support channels selected.',
        'support.chooseChannels': '⌬',
        'support.channelIdPlaceholder': 'Channel ID',
        'support.add': 'Add',
        'support.channelPopupTitle': 'Select support channels',
        'support.channelPopupHint': 'Click a channel to add or remove it. Multiple selections are possible.',
        'support.generalText': 'Plugin activation is managed centrally on the plugins page. This page is for configuration only.',
        'support.noServerGroups': 'No server groups available.',
        'support.unnamedGroup': 'Unnamed group',
        'support.unnamedChannel': 'Unnamed channel',
        'support.validationChannels': 'Select at least one support channel.',
        'support.channelsLoadError': 'Channels could not be loaded.',
        'support.serverGroupsLoadError': 'Server groups could not be loaded.',
        'support.settingsLoadError': 'Failed to load support settings: {error}',
        'support.statusLoadError': 'Failed to load support status: {error}',
        'support.saveError': 'Failed to save support settings: {error}',
        'support.actionDone': 'Support action {action} executed.',
        'support.actionError': 'Support action failed: {error}'
    },
    de: {
        'title.dashboard': 'UC-Framework Dashboard',
        'title.botControl': 'Bot Steuerung',
        'title.counter': 'Counter-Einstellungen',
        'title.plugins': 'Plugins',
        'title.ts3Settings': 'TeamSpeak Einstellungen',
        'title.ts3Connection': 'TS3 Verbindungstest',
        'title.framework': 'Framework',
        'title.discord': 'Discord Einstellungen',
        'title.announcement': 'Ankündigungen',
        'title.support': 'Support-Channels',
        'title.authSettings': 'Auth-Einstellungen',
        'title.afkMover': 'AFK Mover',
        'nav.dashboard': 'Dashboard',
        'nav.teamspeak': 'TeamSpeak',
        'nav.discord': 'Discord',
        'nav.ts3Connection': 'TS3 Verbindung',
        'nav.bot': 'Bot',
        'nav.framework': 'Framework',
        'nav.counter': 'Counter',
        'nav.announcement': 'Ankündigungen',
        'nav.support': 'Support',
        'nav.afkMover': 'AFK Mover',
        'nav.auth': 'Auth',
        'nav.plugins': 'Plugins',
        'nav.settings': 'Einstellungen',
        'nav.logout': 'Logout',
        'nav.language': 'Sprache',
        'nav.languageEnglish': 'Englisch',
        'nav.languageGerman': 'Deutsch',
        'dashboard.heading': 'UC-Dashboard',
        'dashboard.stats': 'Statistiken',
        'dashboard.plugins': 'Plugin-Status',
        'dashboard.settings': 'Einstellungen',
        'dashboard.admins': 'Admins',
        'dashboard.users': 'Benutzer',
        'dashboard.supportersOnline': 'Supporter online',
        'dashboard.membersOnline': 'Mitglieder online',
        'dashboard.loading': 'Lade...',
        'dashboard.statsLoadError': 'Fehler beim Laden der Statistiken.',
        'dashboard.settingsLoadError': 'Fehler beim Laden der Einstellungen.',
        'dashboard.settingsUnavailable': 'Keine Einstellungen verfügbar.',
        'dashboard.ts3Host': 'TS3 Host',
        'dashboard.ts3QueryPort': 'TS3 Query-Port',
        'dashboard.ts3Nickname': 'Bot-Nickname',
        'dashboard.ts3DefaultChannel': 'Standardchannel',
        'dashboard.botStatus': 'Bot-Status',
        'dashboard.botRunning': 'Läuft',
        'dashboard.botStopped': 'Gestoppt',
        'dashboard.botStatusLoadError': 'Bot-Status konnte nicht geladen werden.',
        'plugins.heading': 'Plugins',
        'plugins.description': 'Hier kannst du Plugins aktivieren oder deaktivieren. Die detaillierte Konfiguration bleibt auf den jeweiligen Plugin-Seiten.',
        'plugins.listHeading': 'Plugin-Verwaltung',
        'plugin.adminCounter.name': 'Admin Counter',
        'plugin.adminCounter.description': 'Zählt Online-Admins anhand der konfigurierten Servergruppen.',
        'plugin.memberCounter.name': 'Mitglieder-Counter',
        'plugin.memberCounter.description': 'Zählt Online-Mitglieder und berücksichtigt ausgeschlossene Gruppen oder Nicknames.',
        'plugin.afkMover.name': 'AFK Mover',
        'plugin.afkMover.description': 'Verschiebt inaktive Benutzer nach der konfigurierten Zeit in den AFK-Channel.',
        'plugin.combinedStats.name': 'Kombinierte Statistiken',
        'plugin.combinedStats.description': 'Stellt kombinierte Admin- und Mitgliederstatistiken für Dashboard und API bereit.',
        'plugin.supportControl.name': 'Support-Channel Steuerung',
        'plugin.supportControl.description': 'Öffnet und schließt Support-Channels per Berechtigungen, mit optionalen Servernachrichten sowie Join-Pokes und Auto-Zeitplan.',
        'common.yes': 'Ja',
        'common.no': 'Nein',
        'common.enabled': 'Aktiviert',
        'common.provider': 'Anbieter',
        'common.user': 'Benutzer',
        'common.active': 'aktiv',
        'common.inactive': 'inaktiv',
        'common.saved': 'Gespeichert!',
        'common.error': 'Fehler',
        'common.loading': 'Lade...',
        'common.close': 'Schließen',
        'common.remove': 'Entfernen',
        'common.noEntries': 'Keine Einträge.',
        'common.channelSelectionConnectedHint': 'Channel-Auswahl ist nur möglich, wenn der Bot mit TeamSpeak verbunden ist.',
        'framework.heading': 'Framework',
        'framework.description': 'Hier siehst du Framework-Metadaten, wählst die aktive Kommunikationsplattform und kannst den kompletten Framework-Prozess neu starten.',
        'framework.platformTitle': 'Plattform-Setup',
        'framework.platformHint': 'Wähle, ob das Framework für TeamSpeak, Discord oder beide Plattformen gleichzeitig vorbereitet werden soll.',
        'framework.platformMode': 'Kommunikationsplattformen',
        'framework.platformTeamspeak': 'TeamSpeak',
        'framework.platformDiscord': 'Discord',
        'framework.platformBoth': 'TeamSpeak und Discord',
        'framework.savePlatforms': 'Plattform-Setup speichern',
        'framework.settingsLoaded': 'Framework-Plattformeinstellungen geladen.',
        'framework.settingsLoadError': 'Framework-Einstellungen konnten nicht geladen werden: {error}',
        'framework.saved': 'Framework-Plattformeinstellungen gespeichert.',
        'framework.saveError': 'Framework-Einstellungen konnten nicht gespeichert werden: {error}',
        'framework.infoTitle': 'Framework-Informationen',
        'framework.name': 'Name',
        'framework.version': 'Version',
        'framework.updatedAt': 'Zuletzt aktualisiert',
        'framework.latestVersion': 'Neueste Version',
        'framework.upToDate': 'Aktuell',
        'framework.infoLoaded': 'Framework-Informationen geladen.',
        'framework.infoLoadError': 'Framework-Informationen konnten nicht geladen werden: {error}',
        'framework.restartTitle': 'Framework-Neustart',
        'framework.restartHint': 'Startet den kompletten Framework-Prozess neu, inklusive Web-API und Laufzeitdiensten. (Funktioniert nicht bei alle systeme)',
        'framework.restartButton': 'Framework neu starten',
        'framework.restartConfirm': 'Möchtest du das komplette Framework jetzt wirklich neu starten?',
        'framework.restartScheduled': 'Framework-Neustart wurde eingeplant. Diese Seite wird in Kürze neu geladen.',
        'framework.restartFailed': 'Framework-Neustart fehlgeschlagen.',
        'framework.restartError': 'Framework konnte nicht neu gestartet werden: {error}',
        'announcement.heading': 'Server-Ankündigungen',
        'announcement.description': 'Sende Ankündigungen an alle verbundenen User. Konfiguriere die Ankündigungsnachricht und stelle sie so ein, dass sie einmalig, zu einer festen Zeit oder in regelmäßigen Abständen versendet wird.',
        'announcement.settingsTitle': 'Ankündigungsnachricht',
        'announcement.message': 'Nachricht',
        'announcement.messagePlaceholder': 'Schreibe hier deine Ankündigung...',
        'announcement.messageEmpty': 'Bitte zuerst eine Nachricht eingeben.',
        'announcement.messageTooLong': 'Die Nachricht darf maximal 500 Zeichen lang sein.',
        'announcement.scheduleTitle': 'Zeitplanung',
        'announcement.repeatTitle': 'Ankündigung wiederholen',
        'announcement.repeatHint': 'Versende die Ankündigung automatisch zu einer festen oder wiederkehrenden Zeit. Deaktiviert = einmalig versenden.',
        'announcement.scheduleMode': 'Zeitplan-Modus',
        'announcement.modeOnce': 'Jetzt einmalig senden',
        'announcement.modeInterval': 'Alle X Minuten wiederholen',
        'announcement.modeTime': 'Zu einer bestimmten Uhrzeit versenden',
        'announcement.intervalTime': 'Intervall (Minuten)',
        'announcement.intervalHint': 'Mindestens 10 Minuten',
        'announcement.intervalCount': 'Wie oft versenden',
        'announcement.fixedTime': 'Tägliche Sendezeit',
        'announcement.invalidTime': 'Ungültiges Zeitformat. Verwende HH:MM.',
        'announcement.save': 'Ankündigungseinstellungen speichern',
        'announcement.saved': 'Ankündigungseinstellungen gespeichert.',
        'announcement.loaded': 'Ankündigungseinstellungen geladen.',
        'announcement.manualTitle': 'Jetzt senden',
        'announcement.manualHint': 'Versende die aktuelle Ankündigungsnachricht sofort an alle verbundenen User.',
        'announcement.lastSent': 'Zuletzt gesendet',
        'announcement.neverSent': 'Noch nie gesendet',
        'announcement.sendNow': 'Ankündigung jetzt senden',
        'announcement.sent': 'Ankündigung erfolgreich gesendet.',
        'announcement.refresh': 'Status aktualisieren',
        'announcement.settingsLoadError': 'Ankündigungseinstellungen konnten nicht geladen werden: {error}',
        'announcement.statusLoadError': 'Ankündigungsstatus konnte nicht geladen werden: {error}',
        'announcement.saveError': 'Ankündigungseinstellungen konnten nicht gespeichert werden: {error}',
        'announcement.sendError': 'Ankündigung konnte nicht gesendet werden: {error}',
        'bot.heading': 'Bot Steuerung',
        'bot.description': 'Hier kannst du den Bot starten, stoppen und neu starten. Der Watchdog startet den Bot automatisch neu, wenn er abstürzt oder ein Start fehlschlägt, und versucht dies alle 60-120 Sekunden.',
        'bot.status': 'Bot Status',
        'bot.running': 'Laufend',
        'bot.pid': 'PID',
        'bot.desiredState': 'Gewünschter Zustand',
        'bot.lastAction': 'Letzte Aktion',
        'bot.lastError': 'Letzter Fehler',
        'bot.start': 'Bot starten',
        'bot.stop': 'Bot stoppen',
        'bot.restart': 'Bot neu starten',
        'bot.watchdog': 'Watchdog',
        'bot.watchdogRunning': 'Watchdog läuft',
        'bot.watchdogPid': 'Watchdog PID',
        'bot.startWatchdog': 'Watchdog starten',
        'bot.stopWatchdog': 'Watchdog stoppen',
        'bot.refreshStatus': 'Status aktualisieren',
        'bot.system': 'System',
        'bot.os': 'Betriebssystem',
        'bot.arch': 'Architektur',
        'bot.workingDir': 'Arbeitsverzeichnis',
        'bot.botFile': 'Bot Datei',
        'bot.supervisorScript': 'Supervisor Skript',
        'bot.watchdogScript': 'Watchdog Skript',
        'bot.logFile': 'Log Datei',
        'bot.logs': 'Log Ansicht',
        'bot.lastLines': 'Letzte Zeilen',
        'bot.filter': 'Filter',
        'bot.loadLogs': 'Logs laden',
        'bot.noLogLines': 'Keine Log-Zeilen vorhanden.',
        'bot.logsLoaded': 'Log geladen ({lines} Zeilen, Filter: {level}).',
        'bot.actionExecuted': 'Aktion {action} ausgeführt.',
        'bot.watchdogActionExecuted': 'Watchdog {action} ausgeführt.',
        'bot.statusLoadError': 'Status konnte nicht geladen werden: {error}',
        'bot.logsLoadError': 'Logs konnten nicht geladen werden: {error}',
        'log.level.all': 'Alle',
        'log.level.error': 'Error',
        'log.level.warn': 'Warn',
        'log.level.info': 'Info',
        'log.level.debug': 'Debug',
        'counter.heading': 'Counter Einstellungen',
        'counter.adminCounter': 'Admin Counter',
        'counter.adminDescription': 'Server-Gruppen-IDs, die als Admins gezählt werden.',
        'counter.memberCounter': 'Mitglieder-Counter',
        'counter.memberDescription': 'Alle Benutzer werden gezählt, außer den folgenden Ausnahmen.',
        'counter.excludedGroups': 'Ausgeschlossene Gruppen',
        'counter.excludedGroupsDesc': 'Benutzer mit diesen Server-Gruppen-IDs werden nicht mitgezählt.',
        'counter.excludedNicknames': 'Ausgeschlossene Nicknames',
        'counter.excludedNicknamesDesc': 'Benutzer mit exakt diesen Nicknames werden nicht mitgezählt.',
        'counter.groupId': 'Gruppen-ID',
        'counter.nickname': 'Nickname',
        'counter.renameSection': 'Channel bei Count-Änderung umbenennen',
        'counter.renameHint': 'Nutze einen freien Platzhalter. Beispiel: Template "cspacer7 %%" und Token "%%".',
        'counter.renameChannelId': 'Zu benennende Channel-ID',
        'counter.chooseRenameChannel': '⌬',
        'counter.renameTemplate': 'Neues Channelnamen-Template',
        'counter.renameToken': 'Count-Platzhalter (z. B. {count} oder %%)',
        'counter.channelPopupTitle': 'Verfügbare Channels',
        'counter.channelPopupHint': 'Wähle einen Channel aus. Die gewählte Channel-ID wird direkt gespeichert.',
        'counter.add': 'Hinzufügen',
        'counter.saveAdmin': 'Admin Counter speichern',
        'counter.saveMember': 'Mitglieder-Counter speichern',
        'counter.pluginNotLoaded': 'Plugin ist deaktiviert. Die Aktivierung erfolgt auf der Plugin-Seite.',
        'counter.saveVerifyFailed': 'Server-Antwort enthält deine Rename-Felder nicht. Seite neu laden und erneut speichern.',
        'ts3.heading': 'TeamSpeak Einstellungen',
        'ts3.description': 'Hier kannst du alle TS3-Verbindungs- und Bot-Daten bearbeiten. Nach dem Speichern ist ein Bot-Neustart erforderlich, damit die neue Verbindung aktiv wird.',
        'ts3.server': 'Server',
        'ts3.host': 'TS3 Host-Adresse',
        'ts3.queryPort': 'TS3 Query-Port',
        'ts3.voicePort': 'TS3 Voice-Port',
        'ts3.queryLogin': 'Query Login',
        'ts3.queryUser': 'TS3 Query-Benutzer',
        'ts3.queryPassword': 'TS3 Query-Passwort',
        'ts3.botBehavior': 'Bot Verhalten',
        'ts3.botNickname': 'Bot-Nickname',
        'ts3.defaultChannel': 'Standardchannel',
        'ts3.querySlowmode': 'Query Slowmode (ms)',
        'ts3.chooseChannel': '⌬',
        'ts3.channelPopupTitle': 'Verfügbare Channels',
        'ts3.channelPopupHint': 'Wähle einen Channel als Standardchannel aus.',
        'ts3.channelLoadError': 'Channels konnten nicht geladen werden.',
        'ts3.noChannels': 'Keine Channels verfügbar.',
        'ts3.unnamedChannel': 'Unbenannter Channel',
        'ts3.save': 'TS3 Einstellungen speichern',
        'ts3.loadError': 'TS3 Einstellungen konnten nicht geladen werden.',
        'ts3.loadException': 'Fehler beim Laden der TS3 Einstellungen.',
        'ts3.savedRestart': 'Gespeichert. Bitte Bot neu starten, damit die Änderungen aktiv werden.',
        'ts3.saved': 'Gespeichert.',
        'discord.heading': 'Discord Einstellungen',
        'discord.description': 'Bereite die Discord-Integration vollständig im Webinterface vor. Alle wichtigen Bot-, Server-, Channel- und Rollen-Einstellungen können hier gepflegt werden.',
        'discord.topHelpHtml': '<p><strong>Discord Setup Kurzhilfe</strong></p><p>Bot-Token: geheimer Schluessel, mit dem sich das Framework als dein Discord-Bot anmeldet.</p><p>Application-ID: eindeutige ID deiner Discord-Anwendung (Developer Portal).</p><p>Guild-ID: eindeutige ID des Discord-Servers, auf dem der Bot laufen soll.</p>',
        'discord.enableTitle': 'Discord-Integration',
        'discord.enableHint': 'Aktiviere dies, sobald die Discord-Verbindung im Framework konfiguriert und vorbereitet werden soll.',
        'discord.authTitle': 'Bot-Zugang',
        'discord.authHelpHtml': '<strong>Bot-Zugang</strong><p>Bot-Token: geheimer Schluessel fuer den Bot-Login.</p><p>Application-ID: eindeutige ID der Discord-Anwendung.</p><p>Guild-ID: eindeutige ID deines Discord-Servers.</p>',
        'discord.botToken': 'Bot-Token',
        'discord.botTokenPlaceholder': 'Discord Bot-Token',
        'discord.applicationId': 'Application-ID',
        'discord.applicationIdPlaceholder': 'Discord Application-ID',
        'discord.guildId': 'Guild-ID',
        'discord.guildIdPlaceholder': 'Discord Guild-ID',
        'discord.behaviorTitle': 'Bot-Verhalten',
        'discord.afkTitle': 'Discord AFK-Verhalten',
        'discord.afkKickTitle': 'Inaktive Voice-Nutzer kicken',
        'discord.afkKickHint': 'Discord-spezifisch: Nutzer, die sich für die eingestellte Zeit nicht gemeldet oder ihren Voice-Status geändert haben, werden aus dem Voice getrennt.',
        'discord.afkInactivityMinutes': 'Inaktivitäts-Timeout (Minuten)',
        'discord.botDisplayName': 'Anzeigename',
        'discord.botDisplayNamePlaceholder': 'UC-Framework',
        'discord.statusText': 'Status-Text',
        'discord.statusTextPlaceholder': 'Beobachtet Support-Aktivität',
        'discord.commandPrefix': 'Command-Präfix',
        'discord.commandPrefixPlaceholder': '!',
        'discord.channelsTitle': 'Channels',
        'discord.channelIdPlaceholder': 'Discord Channel-ID',
        'discord.chooseChannel': 'Channel wählen',
        'discord.logChannelId': 'Log-Channel-ID',
        'discord.announcementChannelId': 'Ankündigungs-Channel-ID',
        'discord.supportCategoryId': 'Support-Kategorie-ID',
        'discord.supportLogChannelId': 'Support-Log-Channel-ID',
        'discord.channelPopupTitle': 'Discord-Channel auswählen',
        'discord.channelPopupHint': 'Klicke einen Channel an, um dessen ID für das ausgewählte Feld zu übernehmen.',
        'discord.channelsLoadError': 'Discord-Channels konnten nicht geladen werden.',
        'discord.noChannels': 'Keine Discord-Channels verfügbar.',
        'discord.unnamedChannel': 'Unbenannter Channel',
        'discord.rolesTitle': 'Rollen',
        'discord.adminRoles': 'Admin-Rollen-IDs',
        'discord.supporterRoles': 'Supporter-Rollen-IDs',
        'discord.botRoles': 'Bot-Rollen-IDs',
        'discord.chooseRole': 'Rollen wählen',
        'discord.roleIdPlaceholder': 'Discord Rollen-ID',
        'discord.rolePopupTitle': 'Discord-Rollen auswählen',
        'discord.rolePopupHint': 'Klicke Rollen an, um sie hinzuzufügen oder zu entfernen. Mehrfachauswahl ist möglich.',
        'discord.rolesLoadError': 'Discord-Rollen konnten nicht geladen werden.',
        'discord.noRoles': 'Keine Discord-Rollen verfügbar.',
        'discord.unnamedRole': 'Unbenannte Rolle',
        'discord.addRole': 'Hinzufügen',
        'discord.noRolesSelected': 'Keine Rollen ausgewählt.',
        'discord.save': 'Discord Einstellungen speichern',
        'discord.loaded': 'Discord Einstellungen geladen.',
        'discord.loadError': 'Discord Einstellungen konnten nicht geladen werden: {error}',
        'discord.saved': 'Discord Einstellungen gespeichert.',
        'discord.savedRestart': 'Gespeichert. Bitte Bot neu starten, damit die Änderungen aktiv werden.',
        'discord.saveError': 'Discord Einstellungen konnten nicht gespeichert werden: {error}',
        'ts3conn.heading': 'TS3 Verbindungstest',
        'ts3conn.description': 'Prüfe, ob der Bot aktuell mit TeamSpeak verbunden ist, und starte einen aktiven Verbindungstest.',
        'ts3conn.currentStatus': 'Aktueller Status',
        'ts3conn.connected': 'Verbunden',
        'ts3conn.host': 'Host',
        'ts3conn.port': 'Query-Port',
        'ts3conn.lastCheck': 'Letzte Prüfung',
        'ts3conn.lastError': 'Letzter Fehler',
        'ts3conn.botDisabled': 'Deaktiviert',
        'ts3conn.botDisabledInfo': 'Der Bot ist deaktiviert.',
        'ts3conn.refresh': 'Status aktualisieren',
        'ts3conn.runTest': 'Verbindung testen',
        'ts3conn.statusLoaded': 'Status geladen.',
        'ts3conn.statusLoadError': 'Status konnte nicht geladen werden: {error}',
        'ts3conn.testOk': 'Verbindungstest erfolgreich.',
        'ts3conn.testFailed': 'Verbindungstest fehlgeschlagen.',
        'ts3conn.testError': 'Verbindungstest fehlgeschlagen: {error}',
        'auth.heading': 'Authentifizierung',
        'auth.forcePasswordTitle': 'Passwortwechsel erforderlich',
        'auth.forcePasswordText': 'Beim ersten Login muss ein neues Passwort gesetzt werden. Bis dahin sind andere Bereiche gesperrt.',
        'auth.currentMode': 'Aktueller Modus',
        'auth.loginMode': 'Login-Modus',
        'auth.activeMode': 'Aktiver Login-Modus',
        'auth.modeNone': 'Ohne Login',
        'auth.modeLocal': 'Local-Login',
        'auth.modeRanksystem': 'TSN-Ranksystem-Login',
        'auth.modeLocalRanksystem': 'Local + TSN-Ranksystem-Login',
        'auth.saveMode': 'Login-Modus speichern',
        'auth.modeSaved': 'Login-Modus gespeichert.',
        'auth.noLoginWarningTitle': 'Sicherheitswarnung',
        'auth.noLoginWarningText': 'Login ist deaktiviert. Das ist unsicher und sollte sofort geändert werden.',
        'auth.rankHealth': 'Rank-System-Status',
        'auth.rankHealthText': 'Prüft, ob der konfigurierte Rank-System-Login-Endpoint erreichbar ist.',
        'auth.startHealth': 'Statusprüfung starten',
        'auth.changeLocalPassword': 'Lokales Passwort ändern',
        'auth.currentPassword': 'Aktuelles Passwort',
        'auth.newPassword': 'Neues Passwort',
        'auth.savePassword': 'Passwort speichern',
        'auth.changeLocalUsername': 'Lokalen Benutzernamen ändern',
        'auth.newUsername': 'Neuer Benutzername',
        'auth.saveUsername': 'Benutzernamen speichern',
        'auth.modeLoadError': 'Auth-Modus konnte nicht geladen werden.',
        'auth.healthFailed': 'Statusprüfung fehlgeschlagen.',
        'auth.settingsLoadError': 'Fehler beim Laden der Auth-Einstellungen.',
        'auth.usernameTooShort': 'Neuer Benutzername muss mindestens 3 Zeichen haben.',
        'auth.usernameChanged': 'Benutzername erfolgreich geändert.',
        'auth.passwordTooShort': 'Neues Passwort muss mindestens 8 Zeichen haben.',
        'auth.passwordChanged': 'Passwort erfolgreich geändert.',
        'login.subtitle': 'Webinterface Login',
        'login.text': 'Bitte mit deinen Zugangsdaten anmelden.',
        'login.username': 'Benutzername',
        'login.password': 'Passwort',
        'login.signIn': 'Einloggen',
        'login.tooManyAttempts': 'Zu viele Fehlversuche. Bitte in {seconds} Sekunden erneut versuchen.',
        'login.tooManyAttemptsGeneric': 'Zu viele Fehlversuche. Bitte später erneut versuchen.',
        'login.failed': 'Login fehlgeschlagen. Zugangsdaten prüfen.',
        'login.autoLogout': 'Du wurdest nach 30 Minuten Inaktivität automatisch ausgeloggt.',
        'afk.heading': 'AFK Mover',
        'afk.general': 'Allgemein',
        'afk.generalText': 'Die Aktivierung des Plugins erfolgt zentral auf der Plugin-Seite. Diese Seite dient nur zur Konfiguration.',
        'afk.channelId': 'AFK Channel ID',
        'afk.channelIdPlaceholder': 'Channel ID',
        'afk.chooseTargetChannel': '⌬',
        'afk.targetPreviewLabel': 'Ausgewählter Zielchannel',
        'afk.channelPopupTitle': 'Zielchannel wählen',
        'afk.channelPopupHint': 'Wähle den Channel, in den inaktive User verschoben werden.',
        'afk.chooseExcludedChannel': '⌬',
        'afk.excludedChannelPopupTitle': 'Wähle auszunehmenden Channel',
        'afk.excludedChannelPopupHint': 'Wähle einen Channel aus, der vom AFK-Verschieben ausgenommen werden soll.',
        'afk.channelLoadError': 'Channels konnten nicht geladen werden.',
        'afk.noChannels': 'Keine Channels verfügbar.',
        'afk.unnamedChannel': 'Unbenannter Channel',
        'afk.timeoutMinutes': 'Inaktivitätszeit (Minuten)',
        'afk.returnOnActivityTitle': 'Automatische Rückkehr bei Aktivität',
        'afk.returnOnActivityLabel': 'User bei Aktivität zurück verschieben',
        'afk.returnOnActivityHint': 'Wenn aktiv, werden vom AFK-Mover verschobene User bei neuer Aktivität wieder in ihren vorherigen Channel zurückgesetzt.',
        'afk.excludedChannels': 'Ausgenommene Channels',
        'afk.excludedChannelsText': 'User in diesen Channels werden nicht in den AFK-Channel verschoben.',
        'afk.add': 'Hinzufügen',
        'afk.save': 'Einstellungen speichern',
        'afk.pluginDisabledInfo': 'Plugin ist aktuell deaktiviert. Aktivierung erfolgt auf der Plugin-Seite.',
        'afk.noExceptions': 'Keine Ausnahmen konfiguriert.',
        'afk.loadError': 'Fehler beim Laden der Konfiguration: {error}',
        'support.heading': 'Support-Channel Steuerung',
        'support.description': 'Konfiguriere Support-Channels, Servernachrichten, Join-Pokes, Wartebereich und automatische Öffnungs-/Schliesszeiten. Öffnen/Schliessen wird über Channel-Rechte erzwungen.',
        'support.settingsTitle': 'Support-Einstellungen',
        'support.chooseChannels': '⌬',
        'support.noChannels': 'Keine Channels verfügbar.',
        'support.waitingArea': 'Support-Wartebereich',
        'support.noWaitingArea': 'Kein Wartebereich ausgewählt',
        'support.chooseWaitingArea': '⌬',
        'support.waitingPreviewLabel': 'Ausgewählter Wartebereich',
        'support.waitingPopupTitle': 'Wartebereich wählen',
        'support.waitingPopupHint': 'Wähle genau einen Channel als Wartebereich aus.',
        'support.openPoke': 'Servernachricht beim Öffnen',
        'support.closePoke': 'Servernachricht beim Schliessen',
        'support.joinOpenPoke': 'Poke für User beim Betreten (Support offen)',
        'support.joinClosedPoke': 'Poke für User beim Betreten (Support geschlossen)',
        'support.joinOpenPokePlaceholder': 'Willkommen, der Support ist geöffnet.',
        'support.joinClosedPokePlaceholder': 'Der Support ist aktuell geschlossen.',
        'support.supporterPoke': 'Poke für Supporter wenn ein User beitritt (nur bei offenem Support)',
        'support.supporterPokePlaceholder': 'Neue Support-Anfrage von {user}.',
        'support.supporterGroups': 'Supporter-Servergruppen (Mehrfachauswahl)',
        'support.chooseGroups': '👥',
        'support.groupIdPlaceholder': 'Gruppen-ID',
        'support.groupPopupTitle': 'Supporter-Gruppen auswählen',
        'support.groupPopupHint': 'Klicke auf Gruppen um sie hinzuzufügen oder zu entfernen. Mehrfachauswahl möglich.',
        'support.noGroupsSelected': 'Keine Gruppen ausgewählt',
        'support.openPokePlaceholder': 'Support ist jetzt geöffnet.',
        'support.closePokePlaceholder': 'Support ist jetzt geschlossen.',
        'support.autoTitle': 'Automatisches öffnen und Schliessen',
        'support.autoHint': 'Support-Channels automatisch per Uhrzeit umschalten (lokale Serverzeit).',
        'support.openTime': 'öffnungszeit',
        'support.closeTime': 'Schliesszeit',
        'support.save': 'Support-Einstellungen speichern',
        'support.saved': 'Support-Einstellungen gespeichert.',
        'support.loaded': 'Support-Einstellungen geladen.',
        'support.manualTitle': 'Manuelles öffnen und Schliessen',
        'support.currentState': 'Aktueller Zustand',
        'support.lastAction': 'Letzte Aktion',
        'support.lastError': 'Letzter Fehler',
        'support.openNow': 'Jetzt öffnen',
        'support.closeNow': 'Jetzt schliessen',
        'support.refresh': 'Status aktualisieren',
        'support.stateOpen': 'Geöffnet',
        'support.stateClosed': 'Geschlossen',
        'support.noChannels': 'Keine Channels verfügbar.',
        'support.noChannelsSelected': 'Keine Support-Channels ausgewählt.',
        'support.chooseChannels': '⌬',
        'support.channelIdPlaceholder': 'Channel ID',
        'support.add': 'Hinzufügen',
        'support.channelPopupTitle': 'Support-Channels wählen',
        'support.channelPopupHint': 'Klicke auf einen Channel um ihn hinzuzufügen oder zu entfernen. Mehrfachauswahl möglich.',
        'support.generalText': 'Die Plugin-Aktivierung wird zentral auf der Plugins-Seite verwaltet. Diese Seite dient nur der Konfiguration.',
        'support.noServerGroups': 'Keine Servergruppen verfügbar.',
        'support.unnamedGroup': 'Unbenannte Gruppe',
        'support.unnamedChannel': 'Unbenannter Channel',
        'support.validationChannels': 'Wähle mindestens einen Support-Channel aus.',
        'support.channelsLoadError': 'Channels konnten nicht geladen werden.',
        'support.serverGroupsLoadError': 'Servergruppen konnten nicht geladen werden.',
        'support.settingsLoadError': 'Support-Einstellungen konnten nicht geladen werden: {error}',
        'support.statusLoadError': 'Support-Status konnte nicht geladen werden: {error}',
        'support.saveError': 'Support-Einstellungen konnten nicht gespeichert werden: {error}',
        'support.actionDone': 'Support-Aktion {action} ausgeführt.',
        'support.actionError': 'Support-Aktion fehlgeschlagen: {error}'
        
    }
};

let configuredDefaultLanguage = DEFAULT_LANGUAGE;
let supportedLanguages = ['en', 'de'];
let currentLanguage = normalizeLanguage(localStorage.getItem(LANGUAGE_STORAGE_KEY) || DEFAULT_LANGUAGE);

function normalizeLanguage(language) {
    return language === 'de' ? 'de' : 'en';
}

function isElementNode(value) {
    return value && value.nodeType === Node.ELEMENT_NODE;
}

function hasI18nAttributes(element) {
    return element.hasAttribute('data-i18n') ||
        element.hasAttribute('data-i18n-html') ||
        element.hasAttribute('data-i18n-placeholder') ||
        element.hasAttribute('data-i18n-title');
}

function getTranslationValue(key) {
    const table = TRANSLATIONS[currentLanguage] || TRANSLATIONS[configuredDefaultLanguage] || TRANSLATIONS.en;
    return table[key] ?? TRANSLATIONS.en[key] ?? key;
}

const PLUGIN_I18N_KEYS = {
    AdminCounter: {
        name: 'plugin.adminCounter.name',
        description: 'plugin.adminCounter.description'
    },
    MemberCounter: {
        name: 'plugin.memberCounter.name',
        description: 'plugin.memberCounter.description'
    },
    AfkMover: {
        name: 'plugin.afkMover.name',
        description: 'plugin.afkMover.description'
    },
    CombinedStats: {
        name: 'plugin.combinedStats.name',
        description: 'plugin.combinedStats.description'
    },
    SupportControl: {
        name: 'plugin.supportControl.name',
        description: 'plugin.supportControl.description'
    }
};

function t(key, vars = {}) {
    const template = getTranslationValue(key);
    return template.replace(/\{(\w+)\}/g, (_, name) => {
        if (Object.prototype.hasOwnProperty.call(vars, name)) {
            return String(vars[name]);
        }
        return `{${name}}`;
    });
}

function getPluginLabel(pluginName) {
    const keys = PLUGIN_I18N_KEYS[pluginName];
    if (!keys) {
        return pluginName;
    }

    const value = getTranslationValue(keys.name);
    return value === keys.name ? pluginName : value;
}

function getPluginDescription(plugin) {
    const keys = PLUGIN_I18N_KEYS[plugin.name];
    if (keys) {
        const value = getTranslationValue(keys.description);
        if (value !== keys.description) {
            return value;
        }
    }

    return plugin.description || '';
}

function translateElement(element) {
    const textKey = element.getAttribute('data-i18n');
    const htmlKey = element.getAttribute('data-i18n-html');
    const placeholderKey = element.getAttribute('data-i18n-placeholder');
    const titleKey = element.getAttribute('data-i18n-title');

    if (textKey) {
        element.textContent = t(textKey);
    }
    if (htmlKey) {
        element.innerHTML = t(htmlKey);
    }
    if (placeholderKey) {
        element.setAttribute('placeholder', t(placeholderKey));
    }
    if (titleKey) {
        element.setAttribute('title', t(titleKey));
    }
}

function applyTranslations(root = document) {
    const elements = [];
    if (isElementNode(root) && hasI18nAttributes(root)) {
        elements.push(root);
    }
    if (typeof root.querySelectorAll === 'function') {
        elements.push(...root.querySelectorAll('[data-i18n], [data-i18n-html], [data-i18n-placeholder], [data-i18n-title]'));
    }
    elements.forEach(translateElement);
    document.documentElement.lang = currentLanguage;
    updateLanguageSwitchers();
}

function updateLanguageSwitchers() {
    document.querySelectorAll('[data-language-switcher]').forEach((select) => {
        select.value = currentLanguage;
    });
}

function initLanguageSwitchers(root = document) {
    const selects = typeof root.querySelectorAll === 'function'
        ? root.querySelectorAll('[data-language-switcher]')
        : [];

    selects.forEach((select) => {
        if (select.dataset.languageBound === '1') {
            select.value = currentLanguage;
            return;
        }
        select.dataset.languageBound = '1';
        select.value = currentLanguage;
        select.addEventListener('change', async () => {
            const previousLanguage = currentLanguage;
            const ok = await setLanguage(select.value, true);
            if (!ok) {
                currentLanguage = previousLanguage;
                updateLanguageSwitchers();
                applyTranslations();
            }
        });
    });
}

async function loadLanguageSettings() {
    try {
        const res = await fetch('/api/settings/language');
        if (!res.ok) {
            throw new Error(await res.text());
        }
        const data = await res.json();
        supportedLanguages = Array.isArray(data.supported_languages) && data.supported_languages.length
            ? [...new Set(data.supported_languages.map(normalizeLanguage))]
            : ['en', 'de'];
        configuredDefaultLanguage = normalizeLanguage(data.language || data.default_language || DEFAULT_LANGUAGE);
        currentLanguage = supportedLanguages.includes(configuredDefaultLanguage) ? configuredDefaultLanguage : DEFAULT_LANGUAGE;
        localStorage.setItem(LANGUAGE_STORAGE_KEY, currentLanguage);
    } catch (_) {
        currentLanguage = normalizeLanguage(localStorage.getItem(LANGUAGE_STORAGE_KEY) || configuredDefaultLanguage);
    }

    applyTranslations();
    initLanguageSwitchers();
}

async function setLanguage(language, persist = true) {
    const normalized = normalizeLanguage(language);
    if (!supportedLanguages.includes(normalized)) {
        return false;
    }

    currentLanguage = normalized;
    localStorage.setItem(LANGUAGE_STORAGE_KEY, currentLanguage);
    applyTranslations();

    if (!persist) {
        return true;
    }

    try {
        const res = await fetch('/api/settings/language', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ language: normalized })
        });
        if (!res.ok) {
            throw new Error(await res.text());
        }
        const data = await res.json();
        configuredDefaultLanguage = normalizeLanguage(data.language || data.default_language || normalized);
        supportedLanguages = Array.isArray(data.supported_languages) && data.supported_languages.length
            ? [...new Set(data.supported_languages.map(normalizeLanguage))]
            : supportedLanguages;
        currentLanguage = configuredDefaultLanguage;
        localStorage.setItem(LANGUAGE_STORAGE_KEY, currentLanguage);
        applyTranslations();
        return true;
    } catch (_) {
        return false;
    }
}

function getCurrentLanguage() {
    return currentLanguage;
}

function emitLanguageChanged() {
    window.dispatchEvent(new CustomEvent('uc-language-changed', {
        detail: { language: currentLanguage }
    }));
}

document.addEventListener('DOMContentLoaded', () => {
    applyTranslations();
    initLanguageSwitchers();
    loadLanguageSettings();
});

const _originalSetLanguage = setLanguage;
setLanguage = async function (language, persist = true) {
    const changed = await _originalSetLanguage(language, persist);
    if (changed) {
        emitLanguageChanged();
    }
    return changed;
};

const _originalLoadLanguageSettings = loadLanguageSettings;
loadLanguageSettings = async function () {
    await _originalLoadLanguageSettings();
    emitLanguageChanged();
};

window.t = t;
window.applyTranslations = applyTranslations;
window.initLanguageSwitchers = initLanguageSwitchers;
window.setLanguage = setLanguage;
window.getCurrentLanguage = getCurrentLanguage;
window.getPluginLabel = getPluginLabel;
window.getPluginDescription = getPluginDescription;
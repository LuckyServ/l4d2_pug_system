package database

import (
	"fmt"
	"sync"
	"database/sql"
	"../settings"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var dbConn *sql.DB;
var dbErr error;

type DatabasePlayer struct {
	SteamID64				string
	NicknameBase64			string
	AvatarSmall				string
	AvatarBig				string
	Mmr						int
	MmrUncertainty			float32
	LastGameResult			int
	Access					int //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	ProfValidated			bool
	RulesAccepted			bool
	InitialGames			int
	Twitch					string
	CustomMapsConfirmed		int64
	LastCampaignsPlayed		string
}

type DatabaseSession struct {
	SessionID		string
	SteamID64		string
	Since			int64 //in ms
}

type DatabaseBanRecord struct {
	NicknameBase64		string //permanent nickname
	Access				int
	SteamID64			string
	BannedBySteamID64	string
	CreatedAt			int64 //unix timestamp in milliseconds
	AcceptedAt			int64 //unix timestamp in milliseconds
	BanLength			int64 //unix time in milliseconds
	BanReasonBase64		string
}

type DatabaseGameLog struct {
	ID					string
	Valid				bool
	CreatedAt			int64
	PlayersA			string
	PlayersB			string
	TeamAScores			int
	TeamBScores			int
	ConfoglConfig		string
	CampaignName		string
	ServerIP			string
	Pings				string
}

type DatabaseAnticheatLog struct {
	Index				int			`json:"idx"`
	LogLineBase64		string		`json:"logline"`
}

var MuDatabase sync.RWMutex;


func DatabaseConnect() bool {
	dbConn, dbErr = sql.Open("postgres", "postgres://"+settings.DatabaseUsername+":"+settings.DatabasePassword+"@"+settings.DatabaseHost+":"+settings.DatabasePort+"/"+settings.DatabaseName);
	if (dbErr != nil) {
		fmt.Printf("Error connecting to database (Open): %s\n", dbErr);
		return false;
	}
	dbErr = dbConn.Ping();
	if (dbErr != nil) {
		fmt.Printf("Error connecting to database (Ping): %s\n", dbErr);
		return false;
	}
	fmt.Printf("Database connection successfull\n");
	//LogToFile("Database connection successfull");
	return true;
}

func AddPlayer(oPlayer DatabasePlayer) {
	MuDatabase.Lock();
	//Delete if already registered (shouldn't happen, but just in case)
	dbQueryDelete, errQueryDelete := dbConn.Query("DELETE FROM players_list WHERE steamid64 = '"+oPlayer.SteamID64+"';");
	if (errQueryDelete == nil) {
		dbQueryDelete.Close();
	} else {LogToFile("Error deleting player at AddPlayer: "+oPlayer.SteamID64);};
	//Add player
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO players_list(steamid64, base64nickname, avatar_small, avatar_big, mmr, mmr_uncertainty, last_game_result, access, prof_validated, rules_accepted, initial_games, twitch, custom_maps, map_history) VALUES ('"+oPlayer.SteamID64+"', '"+oPlayer.NicknameBase64+"', '"+oPlayer.AvatarSmall+"', '"+oPlayer.AvatarBig+"', "+fmt.Sprintf("%d", oPlayer.Mmr)+", "+fmt.Sprintf("%.06f", oPlayer.MmrUncertainty)+", "+fmt.Sprintf("%d", oPlayer.LastGameResult)+", "+fmt.Sprintf("%d", oPlayer.Access)+", "+fmt.Sprintf("%v", oPlayer.ProfValidated)+", "+fmt.Sprintf("%v", oPlayer.RulesAccepted)+", 0, '"+oPlayer.Twitch+"', "+fmt.Sprintf("%d", oPlayer.CustomMapsConfirmed)+", '"+oPlayer.LastCampaignsPlayed+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error inserting player at AddPlayer: "+oPlayer.SteamID64);};
	MuDatabase.Unlock();
}

func UpdatePlayer(oPlayer DatabasePlayer) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("UPDATE players_list SET base64nickname = '"+oPlayer.NicknameBase64+"', avatar_small = '"+oPlayer.AvatarSmall+"', avatar_big = '"+oPlayer.AvatarBig+"', mmr = "+fmt.Sprintf("%d", oPlayer.Mmr)+", mmr_uncertainty = "+fmt.Sprintf("%.06f", oPlayer.MmrUncertainty)+", last_game_result = "+fmt.Sprintf("%d", oPlayer.LastGameResult)+", access = "+fmt.Sprintf("%d", oPlayer.Access)+", prof_validated = "+fmt.Sprintf("%v", oPlayer.ProfValidated)+", rules_accepted = "+fmt.Sprintf("%v", oPlayer.RulesAccepted)+", twitch = '"+oPlayer.Twitch+"', custom_maps = "+fmt.Sprintf("%d", oPlayer.CustomMapsConfirmed)+", map_history = '"+oPlayer.LastCampaignsPlayed+"' WHERE steamid64 = '"+oPlayer.SteamID64+"';");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error updating player at UpdatePlayer: "+oPlayer.SteamID64);};
	MuDatabase.Unlock();
}

func UpdateInitialGames(oPlayer DatabasePlayer) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("UPDATE players_list SET initial_games = "+fmt.Sprintf("%d", oPlayer.InitialGames)+" WHERE steamid64 = '"+oPlayer.SteamID64+"';");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error updating player at UpdateInitialGames: "+oPlayer.SteamID64);};
	MuDatabase.Unlock();
}

func AddSession(oSession DatabaseSession) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO sessions_list(session_key, steamid64, since_milli) VALUES ('"+oSession.SessionID+"', '"+oSession.SteamID64+"', "+fmt.Sprintf("%d", oSession.Since)+");");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error adding session at AddSession: "+oSession.SteamID64);};
	MuDatabase.Unlock();
}

func GetMmrShift() int {
	MuDatabase.RLock();
	var iMmrShift int;
	dbQuery, errQuery := dbConn.Query("SELECT the_value FROM mmr_shift LIMIT 1;");
	if (errQuery == nil) {
		for (dbQuery.Next()) {
			dbQuery.Scan(&iMmrShift);
		}
		dbQuery.Close();
	} else {LogToFile("Error getting mmr shift at GetMmrShift");};
	MuDatabase.RUnlock();
	return iMmrShift;
}

func ShiftMmr(iShift int) {
	MuDatabase.Lock();
	dbQueryShift, errDbQueryShift := dbConn.Query("UPDATE mmr_shift SET the_value = the_value + "+fmt.Sprintf("%d", iShift)+";");
	if (errDbQueryShift == nil) {
		dbQueryShift.Close();
	} else {LogToFile("Error updating mmr shift at ShiftMmr");};
	dbQueryPlayers, errDbQueryPlayers := dbConn.Query("UPDATE players_list SET mmr = mmr + "+fmt.Sprintf("%d", iShift)+" WHERE prof_validated = true;");
	if (errDbQueryPlayers == nil) {
		dbQueryPlayers.Close();
	} else {LogToFile("Error updating players list at ShiftMmr");};
	MuDatabase.Unlock();
}

func IncreaseUncertainty() {
	MuDatabase.Lock();
	dbQueryIncrease, errDbQueryIncrease := dbConn.Query("UPDATE players_list SET mmr_uncertainty = mmr_uncertainty + "+fmt.Sprintf("%.06f", settings.IncreaseMmrUncertainty)+";");
	if (errDbQueryIncrease == nil) {
		dbQueryIncrease.Close();
	} else {LogToFile("Error increasing mmr uncertainty at IncreaseUncertainty");};
	dbQueryCut, errDbQueryCut := dbConn.Query("UPDATE players_list SET mmr_uncertainty = "+fmt.Sprintf("%.06f", settings.DefaultMmrUncertainty)+" WHERE mmr_uncertainty > "+fmt.Sprintf("%.06f", settings.DefaultMmrUncertainty)+";");
	if (errDbQueryCut == nil) {
		dbQueryCut.Close();
	} else {LogToFile("Error cutting mmr uncertainty at IncreaseUncertainty");};
	MuDatabase.Unlock();
}

func AddBanRecord(oBanRecord DatabaseBanRecord) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO banlist(steamid64, access, steam_name, banned_by, created_on, accepted_on, banlength, banreason) VALUES ('"+oBanRecord.SteamID64+"', "+fmt.Sprintf("%d", oBanRecord.Access)+", '"+oBanRecord.NicknameBase64+"', '"+oBanRecord.BannedBySteamID64+"', "+fmt.Sprintf("%d", oBanRecord.CreatedAt)+", "+fmt.Sprintf("%d", oBanRecord.AcceptedAt)+", "+fmt.Sprintf("%d", oBanRecord.BanLength)+", '"+oBanRecord.BanReasonBase64+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error adding ban at AddBanRecord: "+oBanRecord.SteamID64);};
	MuDatabase.Unlock();
}

func UpdateBanRecord(oBanRecord DatabaseBanRecord) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("UPDATE banlist SET steamid64 = '"+oBanRecord.SteamID64+"', access = "+fmt.Sprintf("%d", oBanRecord.Access)+", steam_name = '"+oBanRecord.NicknameBase64+"', banned_by = '"+oBanRecord.BannedBySteamID64+"', accepted_on = "+fmt.Sprintf("%d", oBanRecord.AcceptedAt)+", banlength = "+fmt.Sprintf("%d", oBanRecord.BanLength)+", banreason = '"+oBanRecord.BanReasonBase64+"' WHERE created_on = "+fmt.Sprintf("%d", oBanRecord.CreatedAt)+";");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error updating ban at UpdateBanRecord: "+oBanRecord.SteamID64);};
	MuDatabase.Unlock();
}

func DeleteBanRecord(i64CreatedAt int64) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("DELETE FROM banlist WHERE created_on = "+fmt.Sprintf("%d", i64CreatedAt)+";");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error deleting ban at DeleteBanRecord: "+fmt.Sprintf("%d", i64CreatedAt));};
	MuDatabase.Unlock();
}

func RemoveSession(sSessID string) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("DELETE FROM sessions_list WHERE session_key = '"+sSessID+"';");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error removing session at RemoveSession: "+sSessID);};
	MuDatabase.Unlock();
}

func LogGame(oGame DatabaseGameLog) {
	MuDatabase.Lock();
	//Delete if already logged (shouldn't happen, but just in case)
	dbQueryDelete, errQueryDelete := dbConn.Query("DELETE FROM games_log WHERE game_id = '"+oGame.ID+"';");
	if (errQueryDelete == nil) {
		dbQueryDelete.Close();
	} else {LogToFile("Error deleting game log at LogGame: "+oGame.ID);};
	//Add player
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO games_log(game_id, game_valid, created_at, players_a, players_b, team_a_scores, team_b_scores, confogl_config, campaign_name, server_ip, glpings) VALUES ('"+oGame.ID+"', "+fmt.Sprintf("%v", oGame.Valid)+", "+fmt.Sprintf("%d", oGame.CreatedAt)+", '"+oGame.PlayersA+"', '"+oGame.PlayersB+"', "+fmt.Sprintf("%d", oGame.TeamAScores)+", "+fmt.Sprintf("%d", oGame.TeamBScores)+", '"+oGame.ConfoglConfig+"', '"+oGame.CampaignName+"', '"+oGame.ServerIP+"', '"+oGame.Pings+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error inserting game log at LogGame: "+oGame.ID);};
	MuDatabase.Unlock();
}

func AntiCheatLog(sLogLineBase64 string) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO cheat_log(clogline) VALUES ('"+sLogLineBase64+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error inserting anticheat log at AntiCheatLog: "+sLogLineBase64);};
	MuDatabase.Unlock();
}

func GameServerChatLog(i64Time int64, sGameID string, sSteamID64 string, sTextBase64 string) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO gs_chat_log(created_at, gameid, steamid64, logline) VALUES ("+fmt.Sprintf("%d", i64Time)+", '"+sGameID+"', '"+sSteamID64+"', '"+sTextBase64+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error inserting chat log at GameServerChatLog: "+sSteamID64);};
	MuDatabase.Unlock();
}

func PublicChatLog(i64Time int64, sNicknameBase64 string, sSteamID64 string, sTextBase64 string) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO pub_chat_log(created_at, nickname, steamid64, logline) VALUES ("+fmt.Sprintf("%d", i64Time)+", '"+sNicknameBase64+"', '"+sSteamID64+"', '"+sTextBase64+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error inserting chat log at PublicChatLog: "+sSteamID64);};
	MuDatabase.Unlock();
}

func RestorePlayers() []DatabasePlayer {
	MuDatabase.RLock();
	var arDBPlayers []DatabasePlayer;
	dbQueryRetrieve, errQueryRetrieve := dbConn.Query("SELECT steamid64,base64nickname,avatar_small,avatar_big,mmr,mmr_uncertainty,last_game_result,access,prof_validated,rules_accepted,twitch,custom_maps,map_history FROM players_list;");
	if (errQueryRetrieve == nil) {

		for (dbQueryRetrieve.Next()) {
			oDBPlayer := DatabasePlayer{};
			dbQueryRetrieve.Scan(&oDBPlayer.SteamID64, &oDBPlayer.NicknameBase64, &oDBPlayer.AvatarSmall, &oDBPlayer.AvatarBig, &oDBPlayer.Mmr, &oDBPlayer.MmrUncertainty, &oDBPlayer.LastGameResult, &oDBPlayer.Access, &oDBPlayer.ProfValidated, &oDBPlayer.RulesAccepted, &oDBPlayer.Twitch, &oDBPlayer.CustomMapsConfirmed, &oDBPlayer.LastCampaignsPlayed);
			arDBPlayers = append(arDBPlayers, oDBPlayer);
		}

		dbQueryRetrieve.Close();
	}
	MuDatabase.RUnlock();
	return arDBPlayers;
}

func RestoreBans() []DatabaseBanRecord {
	MuDatabase.RLock();
	var arDBBanRecords []DatabaseBanRecord;
	dbQueryRetrieve, errQueryRetrieve := dbConn.Query("SELECT steamid64,access,steam_name,banned_by,created_on,accepted_on,banlength,banreason FROM banlist ORDER BY created_on;"); //ordering is important
	if (errQueryRetrieve == nil) {

		for (dbQueryRetrieve.Next()) {
			oDBBanRecord := DatabaseBanRecord{};
			dbQueryRetrieve.Scan(&oDBBanRecord.SteamID64, &oDBBanRecord.Access, &oDBBanRecord.NicknameBase64, &oDBBanRecord.BannedBySteamID64, &oDBBanRecord.CreatedAt, &oDBBanRecord.AcceptedAt, &oDBBanRecord.BanLength, &oDBBanRecord.BanReasonBase64);
			arDBBanRecords = append(arDBBanRecords, oDBBanRecord);
		}

		dbQueryRetrieve.Close();
	}
	MuDatabase.RUnlock();
	return arDBBanRecords;
}

func RestoreSessions() []DatabaseSession {
	MuDatabase.RLock();
	var arDBSessions []DatabaseSession;
	dbQueryRetrieve, errQueryRetrieve := dbConn.Query("SELECT session_key,steamid64,since_milli FROM sessions_list;");
	if (errQueryRetrieve == nil) {

		for (dbQueryRetrieve.Next()) {
			oDBSession := DatabaseSession{};
			dbQueryRetrieve.Scan(&oDBSession.SessionID, &oDBSession.SteamID64, &oDBSession.Since);
			arDBSessions = append(arDBSessions, oDBSession);
		}

		dbQueryRetrieve.Close();
	}
	MuDatabase.RUnlock();
	return arDBSessions;
}

func GetAnticheatLogs() []DatabaseAnticheatLog {
	MuDatabase.RLock();
	var arDBAnticheatLogs []DatabaseAnticheatLog;
	dbQueryRetrieve, errQueryRetrieve := dbConn.Query("SELECT cidx,clogline FROM cheat_log ORDER BY cidx DESC LIMIT 10000;");
	if (errQueryRetrieve == nil) {

		for (dbQueryRetrieve.Next()) {
			oDBAnticheatLog := DatabaseAnticheatLog{};
			dbQueryRetrieve.Scan(&oDBAnticheatLog.Index, &oDBAnticheatLog.LogLineBase64);
			arDBAnticheatLogs = append(arDBAnticheatLogs, oDBAnticheatLog);
		}

		dbQueryRetrieve.Close();
	}
	MuDatabase.RUnlock();
	return arDBAnticheatLogs;
}

func LogToFile(sText string) {
	f, err := os.OpenFile(settings.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777);
	if (err == nil) {
		logger := log.New(f, "db ", log.LstdFlags);
		logger.Println(sText);
		f.Close();
	}
}

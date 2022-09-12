package database

import (
	"fmt"
	"database/sql"
	"../settings"
	_ "github.com/lib/pq"
)

var dbConn *sql.DB;
var dbErr error;

type DatabasePlayer struct {
	SteamID64		string
	NicknameBase64	string
	Mmr				int
	MmrUncertainty	float32
	Access			int //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	ProfValidated	bool
	RulesAccepted	bool
}

type DatabaseSession struct {
	SessionID		string
	SteamID64		string
	Since			int64 //in ms
}

type DatabaseBanRecord struct {
	NicknameBase64		string //permanent nickname
	SteamID64			string
	BannedBySteamID64	string
	CreatedAt			int64 //unix timestamp in milliseconds
	AcceptedAt			int64 //unix timestamp in milliseconds
	BanLength			int64 //unix time in milliseconds
	BanReasonBase64		string
}


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
	return true;
}

func AddPlayer(oPlayer DatabasePlayer) {
	//Delete if already registered (shouldn't happen)
	dbQueryDelete, errQueryDelete := dbConn.Query("DELETE FROM players_list WHERE steamid64 = '"+oPlayer.SteamID64+"';");
	if (errQueryDelete == nil) {
		dbQueryDelete.Close();
	}
	//Add player
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO players_list(steamid64, base64nickname, mmr, mmr_uncertainty, access, prof_validated, rules_accepted) VALUES ("+oPlayer.SteamID64+", '"+oPlayer.NicknameBase64+"', "+fmt.Sprintf("%d", oPlayer.Mmr)+", "+fmt.Sprintf("%.06f", oPlayer.MmrUncertainty)+", "+fmt.Sprintf("%d", oPlayer.Access)+", "+fmt.Sprintf("%v", oPlayer.ProfValidated)+", "+fmt.Sprintf("%v", oPlayer.RulesAccepted)+");");
	if (errDbQuery == nil) {
		dbQuery.Close();
	}
}

func UpdatePlayer(oPlayer DatabasePlayer) {
	dbQuery, errDbQuery := dbConn.Query("UPDATE players_list SET base64nickname = '"+oPlayer.NicknameBase64+"', mmr = "+fmt.Sprintf("%d", oPlayer.Mmr)+", mmr_uncertainty = "+fmt.Sprintf("%.06f", oPlayer.MmrUncertainty)+", access = "+fmt.Sprintf("%d", oPlayer.Access)+", prof_validated = "+fmt.Sprintf("%v", oPlayer.ProfValidated)+", rules_accepted = "+fmt.Sprintf("%v", oPlayer.RulesAccepted)+" WHERE steamid64 = '"+oPlayer.SteamID64+"';");
	if (errDbQuery == nil) {
		dbQuery.Close();
	}
}

func AddSession(oSession DatabaseSession) {
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO sessions_list(session_key, steamid64, since_milli) VALUES ('"+oSession.SessionID+"', '"+oSession.SteamID64+"', "+fmt.Sprintf("%d", oSession.Since)+");");
	if (errDbQuery == nil) {
		dbQuery.Close();
	}
}

func AddBanRecord(oBanRecord DatabaseBanRecord) {
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO banlist(steamid64, steam_name, banned_by, created_on, accepted_on, banlength, banreason) VALUES ('"+oBanRecord.SteamID64+"', '"+oBanRecord.NicknameBase64+"', '"+oBanRecord.BannedBySteamID64+"', "+fmt.Sprintf("%d", oBanRecord.CreatedAt)+", "+fmt.Sprintf("%d", oBanRecord.AcceptedAt)+", "+fmt.Sprintf("%d", oBanRecord.BanLength)+", '"+oBanRecord.BanReasonBase64+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	}
}

func UpdateBanRecord(oBanRecord DatabaseBanRecord) {
	dbQuery, errDbQuery := dbConn.Query("UPDATE banlist SET steamid64 = '"+oBanRecord.SteamID64+"', steam_name = '"+oBanRecord.NicknameBase64+"', banned_by = '"+oBanRecord.BannedBySteamID64+"', accepted_on = "+fmt.Sprintf("%d", oBanRecord.AcceptedAt)+", banlength = "+fmt.Sprintf("%d", oBanRecord.BanLength)+", banreason = '"+oBanRecord.BanReasonBase64+"' WHERE created_on = "+fmt.Sprintf("%d", oBanRecord.CreatedAt)+";");
	if (errDbQuery == nil) {
		dbQuery.Close();
	}
}

func RemoveSession(sSessID string) {
	dbQuery, errDbQuery := dbConn.Query("DELETE FROM sessions_list WHERE session_key = '"+sSessID+"';");
	if (errDbQuery == nil) {
		dbQuery.Close();
	}
}

func RestorePlayers() []DatabasePlayer {
	var arDBPlayers []DatabasePlayer;
	dbQueryRetrieve, errQueryRetrieve := dbConn.Query("SELECT steamid64,base64nickname,mmr,mmr_uncertainty,access,prof_validated,rules_accepted FROM players_list;");
	if (errQueryRetrieve == nil) {

		for (dbQueryRetrieve.Next()) {
			oDBPlayer := DatabasePlayer{};
			dbQueryRetrieve.Scan(&oDBPlayer.SteamID64, &oDBPlayer.NicknameBase64, &oDBPlayer.Mmr, &oDBPlayer.MmrUncertainty, &oDBPlayer.Access, &oDBPlayer.ProfValidated, &oDBPlayer.RulesAccepted);
			arDBPlayers = append(arDBPlayers, oDBPlayer);
		}

		dbQueryRetrieve.Close();
	}
	return arDBPlayers;
}

func RestoreSessions() []DatabaseSession {
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
	return arDBSessions;
}
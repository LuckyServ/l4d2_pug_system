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
	dbQueryAdd, errQueryAdd := dbConn.Query("INSERT INTO players_list(steamid64, base64nickname, mmr, mmr_uncertainty, access, prof_validated, rules_accepted) VALUES ("+oPlayer.SteamID64+", '"+oPlayer.NicknameBase64+"', "+fmt.Sprintf("%d", oPlayer.Mmr)+", "+fmt.Sprintf("%.06f", oPlayer.MmrUncertainty)+", "+fmt.Sprintf("%d", oPlayer.Access)+", "+fmt.Sprintf("%v", oPlayer.ProfValidated)+", "+fmt.Sprintf("%v", oPlayer.RulesAccepted)+");");
	if (errQueryAdd == nil) {
		dbQueryAdd.Close();
	}
}

func UpdatePlayer(oPlayer DatabasePlayer) {
	dbQueryUpdate, errQueryUpdate := dbConn.Query("UPDATE players_list SET base64nickname = '"+oPlayer.NicknameBase64+"', mmr = "+fmt.Sprintf("%d", oPlayer.Mmr)+", mmr_uncertainty = "+fmt.Sprintf("%.06f", oPlayer.MmrUncertainty)+", access = "+fmt.Sprintf("%d", oPlayer.Access)+", prof_validated = "+fmt.Sprintf("%v", oPlayer.ProfValidated)+", rules_accepted = "+fmt.Sprintf("%v", oPlayer.RulesAccepted)+" WHERE steamid64 = '"+oPlayer.SteamID64+"';");
	if (errQueryUpdate == nil) {
		dbQueryUpdate.Close();
	}
}

func AddSession(oSession DatabaseSession) {
	dbQueryAdd, errQueryAdd := dbConn.Query("INSERT INTO sessions_list(session_key, steamid64, since_milli) VALUES ('"+oSession.SessionID+"', '"+oSession.SteamID64+"', "+fmt.Sprintf("%d", oSession.Since)+");");
	if (errQueryAdd == nil) {
		dbQueryAdd.Close();
	}
}

func RemoveSession(sSessID string) {
	dbQueryAdd, errQueryAdd := dbConn.Query("DELETE FROM sessions_list WHERE session_key = '"+sSessID+"';");
	if (errQueryAdd == nil) {
		dbQueryAdd.Close();
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
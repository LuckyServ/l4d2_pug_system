package database

import (
	"fmt"
)

type DatabaseTicket struct {
	TicketID		string		`json:"id"`
	TicketType		int			`json:"type"`
	CreatedBy		string		`json:"created_by"`
	CreatedAt		int64		`json:"created_at"` //milliseconds
	IsClosed		bool		`json:"is_closed"`
}

type DatabaseTicketMessage struct {
	TicketID		string		`json:"ticket_id"`
	MessageBy		string		`json:"by"`
	MessageAt		int64		`json:"at"` //milliseconds
	MessageBase64	string		`json:"text"`
}

const ( //game states
	TicketTypeDummy int = iota
	TicketTypeBehaviour
	TicketTypeCheat
	TicketTypeProtest
	TicketTypeOther
	TicketTypeDummy2
)


func CreateTicket(oTicket DatabaseTicket) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO tickets_list(ticket_id, ticket_type, created_by, cteated_at, is_closed) VALUES ('"+oTicket.TicketID+"', "+fmt.Sprintf("%d", oTicket.TicketType)+", '"+oTicket.CreatedBy+"', "+fmt.Sprintf("%d", oTicket.CreatedAt)+", false);");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error inserting ticket at CreateTicket: "+oTicket.CreatedBy);};
	MuDatabase.Unlock();
}

func CreateMessage(oMessage DatabaseTicketMessage) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("INSERT INTO ticket_messages(ticket_id, message_by, message_at, message_text) VALUES ('"+oMessage.TicketID+"', '"+oMessage.MessageBy+"', "+fmt.Sprintf("%d", oMessage.MessageAt)+", '"+oMessage.MessageBase64+"');");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error inserting ticket message at CreateMessage: "+oMessage.MessageBy);};
	MuDatabase.Unlock();
}

func CloseTicket(sTicketID string) {
	MuDatabase.Lock();
	dbQuery, errDbQuery := dbConn.Query("UPDATE tickets_list SET is_closed = true WHERE ticket_id = '"+sTicketID+"';");
	if (errDbQuery == nil) {
		dbQuery.Close();
	} else {LogToFile("Error updating ticket at CloseTicket: "+sTicketID);};
	MuDatabase.Unlock();
}

func GetOpenedTicketsOfPlayer(sSteamID64 string) []DatabaseTicket {
	MuDatabase.RLock();
	var arDBTickets []DatabaseTicket;
	dbQuery, errDbQuery := dbConn.Query("SELECT ticket_id,ticket_type,cteated_at FROM tickets_list WHERE created_by = '"+sSteamID64+"' AND is_closed = false ORDER BY cteated_at DESC LIMIT 5;");
	if (errDbQuery == nil) {

		for (dbQuery.Next()) {
			oDBTicket := DatabaseTicket{
				CreatedBy:		sSteamID64,
				IsClosed:		false,
			};
			dbQuery.Scan(&oDBTicket.TicketID, &oDBTicket.TicketType, &oDBTicket.CreatedAt);
			arDBTickets = append(arDBTickets, oDBTicket);
		}

		dbQuery.Close();
	}
	MuDatabase.RUnlock();
	return arDBTickets;
}

func GetClosedTicketsOfPlayer(sSteamID64 string) []DatabaseTicket {
	MuDatabase.RLock();
	var arDBTickets []DatabaseTicket;
	dbQuery, errDbQuery := dbConn.Query("SELECT ticket_id,ticket_type,cteated_at FROM tickets_list WHERE created_by = '"+sSteamID64+"' AND is_closed = true ORDER BY cteated_at DESC LIMIT 5;");
	if (errDbQuery == nil) {

		for (dbQuery.Next()) {
			oDBTicket := DatabaseTicket{
				CreatedBy:		sSteamID64,
				IsClosed:		true,
			};
			dbQuery.Scan(&oDBTicket.TicketID, &oDBTicket.TicketType, &oDBTicket.CreatedAt);
			arDBTickets = append(arDBTickets, oDBTicket);
		}

		dbQuery.Close();
	}
	MuDatabase.RUnlock();
	return arDBTickets;
}

func GetAdminTickets(iAccess int) []DatabaseTicket {
	MuDatabase.RLock();
	var arDBTickets []DatabaseTicket;

	var sTicketFilter string;
	if (iAccess == 1) {
		sTicketFilter = fmt.Sprintf(" WHERE ticket_type = %d AND is_closed = false", TicketTypeBehaviour);
	} else if (iAccess == 2) {
		sTicketFilter = fmt.Sprintf(" WHERE ticket_type = %d AND is_closed = false", TicketTypeCheat);
	} else if (iAccess == 3) {
		sTicketFilter = fmt.Sprintf(" WHERE (ticket_type = %d OR ticket_type = %d) AND is_closed = false", TicketTypeBehaviour, TicketTypeCheat);
	} else { //already proved that iAccess > 0
		sTicketFilter = " WHERE is_closed = false";
	}

	dbQuery, errDbQuery := dbConn.Query("SELECT ticket_id,ticket_type,created_by,cteated_at FROM tickets_list"+sTicketFilter+" ORDER BY cteated_at DESC;");
	if (errDbQuery == nil) {

		for (dbQuery.Next()) {
			oDBTicket := DatabaseTicket{
				IsClosed:		false,
			};
			dbQuery.Scan(&oDBTicket.TicketID, &oDBTicket.TicketType, &oDBTicket.CreatedBy, &oDBTicket.CreatedAt);
			arDBTickets = append(arDBTickets, oDBTicket);
		}

		dbQuery.Close();
	}
	MuDatabase.RUnlock();
	return arDBTickets;
}

func GetMessagesOfTicket(sTicketID string) []DatabaseTicketMessage {
	MuDatabase.RLock();
	var arDBTicketMsgs []DatabaseTicketMessage;
	dbQuery, errDbQuery := dbConn.Query("SELECT message_by,message_at,message_text FROM ticket_messages WHERE ticket_id = '"+sTicketID+"' ORDER BY message_at DESC;");
	if (errDbQuery == nil) {

		for (dbQuery.Next()) {
			oDBMsg := DatabaseTicketMessage{
				TicketID:		sTicketID,
			};
			dbQuery.Scan(&oDBMsg.MessageBy, &oDBMsg.MessageAt, &oDBMsg.MessageBase64);
			arDBTicketMsgs = append(arDBTicketMsgs, oDBMsg);
		}

		dbQuery.Close();
	}
	MuDatabase.RUnlock();
	return arDBTicketMsgs;
}

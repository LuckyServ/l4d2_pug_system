forward void L4D2C_GameInfoReceived();
forward void L4D2C_OnGameEnded();
native int L4D2C_GetServerReservation(); //-2 connection error, -1 - havent yet checked, 0 - not reserved, 1 - reserved
native bool L4D2C_IsPlayerGameParticipant(int client);
native bool L4D2C_IsMidRound();


public SharedPlugin __pl_l4d2center =
{
	name = "l4d2center",
	file = "l4d2center.smx",
#if defined REQUIRE_PLUGIN
	required = 1,
#else
	required = 0,
#endif
};

#if !defined REQUIRE_PLUGIN
public void __pl_l4d2center_SetNTVOptional()
{
	MarkNativeAsOptional("L4D2C_GetServerReservation");
	MarkNativeAsOptional("L4D2C_IsPlayerGameParticipant");
	MarkNativeAsOptional("L4D2C_IsMidRound");
}
#endif

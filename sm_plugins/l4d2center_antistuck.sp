#pragma semicolon 1
#pragma newdecls required

#include <sourcemod>
#include <sdktools>
#include <pause>
#undef REQUIRE_PLUGIN
#include <l4d2center>

float Ground_Velocity[3] = {0.0, 0.0, 0.0};
float fDownToFloor[3] = {90.0, 0.0, 0.0};
float fAbleToMove[MAXPLAYERS + 1];
float fNow;
float fMaxTeleportRadius = 200.0;
bool bGlobalAllowedTeleporting;
float fTickrate;
float fMoveCheck;

bool bL4D2CenterAvailable;

public void OnPluginStart() {
	CreateTimer(1.0, Teleport_Callback, 0, TIMER_REPEAT);
	HookEvent("ability_use", Event_AbilityUse);
	HookEvent("round_end", Event_RoundEnd);
	HookEvent("round_start", Event_RoundStart);
	fTickrate = 1.0 / GetTickInterval();
	if (fTickrate >= 50) {
		fMoveCheck = 15.0;
	} else {
		fMoveCheck = 20.0;
	}
}

public void OnAllPluginsLoaded() {
    bL4D2CenterAvailable = LibraryExists("l4d2center");
}

public void OnLibraryRemoved(const char[] sPluginName) {
	if (strcmp(sPluginName, "l4d2center") == 0) {
		bL4D2CenterAvailable = false;
	}
}

public void OnLibraryAdded(const char[] sPluginName) {
	if (strcmp(sPluginName, "l4d2center") == 0) {
		bL4D2CenterAvailable = true;
	}
}

public void Event_RoundEnd(Handle event, const char[] name, bool dontBroadcast) {
	bGlobalAllowedTeleporting = false;
}

public void Event_RoundStart(Handle event, const char[] name, bool dontBroadcast) {
	CreateTimer(5.0, RoundStart_Timer);
}

public Action RoundStart_Timer(Handle timer) {
	bGlobalAllowedTeleporting = true;
	return Plugin_Handled;
}

public void Event_AbilityUse(Handle event, const char[] name, bool dontBroadcast) {
	int client = GetClientOfUserId(GetEventInt(event, "userid"));
	if (client > 0 && IsClientInGame(client) && GetClientTeam(client) > 0) {
		char sAbility[128];
		GetEventString(event, "ability", sAbility, sizeof(sAbility));
		if (StrEqual(sAbility, "ability_vomit")) {
			fAbleToMove[client] = GetEngineTime() + 1.5;
		} else if (StrEqual(sAbility, "ability_throw")) {
			fAbleToMove[client] = GetEngineTime() + 3.5;
		} else if (StrEqual(sAbility, "ability_toungue")) {
			fAbleToMove[client] = GetEngineTime() + 1.0;
		} else if (StrEqual(sAbility, "ability_spit")) {
			fAbleToMove[client] = GetEngineTime() + 1.5;
		}
	}
}

public Action Teleport_Callback(Handle timer, any sheo) {
	if (bL4D2CenterAvailable && L4D2C_GetServerReservation() == 1 && bGlobalAllowedTeleporting && !IsInPause()) {
		for (int i = 1; i <= MaxClients; i++) {
			if (IsClientInGame(i) && GetClientTeam(i) > 1 && IsSafeToTeleport(i)) {
				if (IsEntityStuck(i)) {
					CheckIfPlayerCanMove(i, 0, fMoveCheck, 0.0, 0.0);
				}
			}
		}
	}
	return Plugin_Handled;
}

void CheckIfPlayerCanMove(int iClient, int testID, float X=0.0, float Y=0.0, float Z=0.0) {
	float vecVelo[3];
	float vecOrigin[3];
	GetClientAbsOrigin(iClient, vecOrigin);
	GetEntPropVector(iClient, Prop_Data, "m_vecBaseVelocity", vecVelo);
	vecVelo[0] = vecVelo[0] + X;
	vecVelo[1] = vecVelo[1] + Y;
	vecVelo[2] = vecVelo[2] + Z;
	SetEntPropVector(iClient, Prop_Data, "m_vecBaseVelocity", vecVelo);
	Handle hData = CreateDataPack();
	CreateTimer(0.1, TimerWait, hData);
	WritePackCell(hData, GetClientUserId(iClient));
	WritePackCell(hData, testID);
	WritePackFloat(hData, vecOrigin[0]);
	WritePackFloat(hData, vecOrigin[1]);
	WritePackFloat(hData, vecOrigin[2]);
}

public Action TimerWait(Handle timer, any hData) {
	float vecOrigin[3];
	float vecOriginAfter[3];
	ResetPack(hData, false);
	int iClient = GetClientOfUserId(ReadPackCell(hData));
	if (bGlobalAllowedTeleporting && !IsInPause() && iClient > 0 && IsClientInGame(iClient) && GetClientTeam(iClient) > 0 && IsSafeToTeleport(iClient)) {
		int testID = ReadPackCell(hData);
		vecOrigin[0] = ReadPackFloat(hData);
		vecOrigin[1] = ReadPackFloat(hData);
		vecOrigin[2] = ReadPackFloat(hData);
		GetClientAbsOrigin(iClient, vecOriginAfter);
		if (GetVectorDistance(vecOrigin, vecOriginAfter, false) == 0.0) {
			if(testID == 0) {
				CheckIfPlayerCanMove(iClient, 1, 0.0, 0.0, -1.0 * fMoveCheck);
			} else if(testID == 1) {
				CheckIfPlayerCanMove(iClient, 2, -1.0 * fMoveCheck, 0.0, 0.0);
			} else if(testID == 2) {
				CheckIfPlayerCanMove(iClient, 3, 0.0, fMoveCheck, 0.0);
			} else if(testID == 3) {
				CheckIfPlayerCanMove(iClient, 4, 0.0, -1.0 * fMoveCheck, 0.0);
			} else if(testID == 4) {
				CheckIfPlayerCanMove(iClient, 5, 0.0, 0.0, fMoveCheck);
			} else {
				FixPlayerPosition(iClient);
			}
		}
	}
	CloseHandle(hData);
	return Plugin_Continue;
}

void FixPlayerPosition(int iClient) {
	float pos_Z = -50.0;
	float fRadius = 0.0;
	while (pos_Z <= fMaxTeleportRadius && !TryFixPosition(iClient, fRadius, pos_Z)) {
		fRadius = fRadius + 2.0;
		pos_Z = pos_Z + 2.0;
	}
}

bool TryFixPosition(int iClient, float Radius, float pos_Z) {
	float DegreeAngle;
	float vecPosition[3];
	float vecOrigin[3];
	float vecAngle[3];
	GetClientAbsOrigin(iClient, vecOrigin);
	GetClientEyeAngles(iClient, vecAngle);
	vecPosition[2] = vecOrigin[2] + pos_Z;
	DegreeAngle = -180.0;
	while (DegreeAngle < 180.0) {
		vecPosition[0] = vecOrigin[0] + Radius * Cosine(DegreeAngle * FLOAT_PI / 180.0);
		vecPosition[1] = vecOrigin[1] + Radius * Sine(DegreeAngle * FLOAT_PI / 180.0);
		
		TeleportEntity(iClient, vecPosition, vecAngle, Ground_Velocity);
		if (!IsEntityStuck(iClient) && GetDistanceToFloor(iClient) <= 240.0) {
			return true;
		}
		DegreeAngle += 10.0;
	}
	TeleportEntity(iClient, vecOrigin, vecAngle, Ground_Velocity);
	return false;
}

bool IsEntityStuck(int iEnt) {
	float vecMin[3], vecMax[3], vecOrigin[3];
	GetEntPropVector(iEnt, Prop_Send, "m_vecMins", vecMin);
	GetEntPropVector(iEnt, Prop_Send, "m_vecMaxs", vecMax);
	GetEntPropVector(iEnt, Prop_Send, "m_vecOrigin", vecOrigin);
	Handle hTrace = TR_TraceHullFilterEx(vecOrigin, vecOrigin, vecMin, vecMax, MASK_PLAYERSOLID, TraceEntityFilterSolid);
	bool bTrue = TR_DidHit(hTrace);
	CloseHandle(hTrace);
	return bTrue;
}

float GetDistanceToFloor(int client) {
	float vOrigin[3];
	GetClientEyePosition(client, vOrigin);
	Handle hTrace = TR_TraceRayFilterEx(vOrigin, fDownToFloor, MASK_PLAYERSOLID, RayType_Infinite, TraceEntityFilterSolid);
	if (TR_DidHit(hTrace)) {
		float vFloorPoint[3];
		TR_GetEndPosition(vFloorPoint, hTrace);
		CloseHandle(hTrace);
		return (vOrigin[2] - vFloorPoint[2]);
	}
	CloseHandle(hTrace);
	return 999999.0;
}

public bool TraceEntityFilterSolid(int entity, int contentsMask) {
	if (entity > 0 && entity <= MaxClients) {
		return false;
	}
	int iCollisionType;
	if (entity >= 0 && IsValidEdict(entity) && IsValidEntity(entity)) {
		iCollisionType = GetEntProp(entity, Prop_Send, "m_CollisionGroup");
	}
	if (iCollisionType == 1 || iCollisionType == 11 || iCollisionType == 5) {
		return false;
	}
	return true;
}

bool IsSafeToTeleport(int client) {
	fNow = GetEngineTime();
	if (!IsPlayerAlive(client)) {
		return false;
	} else if (GetEntityMoveType(client) == MOVETYPE_LADDER) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_jockeyAttacker") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_jockeyVictim") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_pounceAttacker") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_pounceVictim") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_carryAttacker") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_carryVictim") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_pummelAttacker") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_pummelVictim") > 0) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_isHangingFromLedge") == 1) {
		return false;
	} else if (GetEntProp(client, Prop_Send, "m_isIncapacitated") == 1) {
		return false;
	} else if (fAbleToMove[client] != 0.0 && fNow < fAbleToMove[client]) {
		return false;
	} else {
		return true;
	}
}

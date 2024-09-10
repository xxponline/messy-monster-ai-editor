package db

import (
	"database/sql"
	"github.com/google/uuid"
	"messy-monster-ai-editor/common"
	"sync"
)

type SqliteDataBase struct {
	sqliteDb *sql.DB

	solutionMgrLocker *sync.RWMutex
	assetSetMgrLocker *sync.RWMutex
	assetMgrLocker    *sync.RWMutex

	assetDocMapLocker sync.Mutex
	assetDocLockerMap map[string]*sync.RWMutex
}

func (db *SqliteDataBase) Initialize(dataSource string) (common.ErrorCode, string) {
	db.solutionMgrLocker = &sync.RWMutex{}
	db.assetMgrLocker = &sync.RWMutex{}
	db.assetSetMgrLocker = &sync.RWMutex{}
	db.assetMgrLocker = &sync.RWMutex{}
	db.assetDocLockerMap = make(map[string]*sync.RWMutex)

	var err error
	db.sqliteDb, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		return common.DataBaseError, err.Error()
	}
	return common.Success, ""
}

func (db *SqliteDataBase) GetSolutionManager(WriteLock bool) (common.ErrorCode, string, ISolutionManager) {

	if WriteLock {
		db.solutionMgrLocker.Lock()
	} else {
		db.solutionMgrLocker.RLock()
	}
	return 0, "", &SqliteSolutionManager{db.sqliteDb, WriteLock, db.solutionMgrLocker}
}

func (db *SqliteDataBase) GetAssetSetManager(WriteLock bool) (common.ErrorCode, string, IAssetSetManager) {
	if WriteLock {
		db.assetSetMgrLocker.Lock()
	} else {
		db.assetSetMgrLocker.RLock()
	}
	return 0, "", &SqliteAssetSetManager{db.sqliteDb, WriteLock, db.assetSetMgrLocker}
}

func (db *SqliteDataBase) GetAssetManager(WriteLock bool) (common.ErrorCode, string, IAssetManager) {
	if WriteLock {
		db.assetMgrLocker.Lock()
	} else {
		db.assetMgrLocker.RLock()
	}
	return 0, "", &SqliteAssetManager{db.sqliteDb, WriteLock, db.assetMgrLocker}
}

func (db *SqliteDataBase) GetAssetDocument(assetId string, WriteLock bool) (common.ErrorCode, string, IAssetDocument) {
	docLocker, ok := db.assetDocLockerMap[assetId]
	if !ok {
		db.assetDocMapLocker.Lock()
		docLocker = &sync.RWMutex{}
		db.assetDocLockerMap[assetId] = docLocker
		db.assetDocMapLocker.Unlock()
	}

	if WriteLock {
		docLocker.Lock()
	} else {
		docLocker.RLock()
	}
	return 0, "", &SqliteAssetDocument{assetId, db.sqliteDb, WriteLock, docLocker}
}

//SolutionManager

type SqliteSolutionManager struct {
	sqliteDb    *sql.DB
	isWriteable bool
	locker      *sync.RWMutex
}

func (solutionMgr *SqliteSolutionManager) Release() {
	if solutionMgr != nil {
		if solutionMgr.isWriteable {
			solutionMgr.locker.Unlock()
		} else {
			solutionMgr.locker.RUnlock()
		}
	}
}

func (solutionMgr *SqliteSolutionManager) ListSolutions() (common.ErrorCode, string, []common.SolutionInfoItem) {
	querySQL := `SELECT id, solutionName FROM ai_solutions`
	{
		statement, err := solutionMgr.sqliteDb.Prepare(querySQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error(), nil
		}
		defer statement.Close()

		rows, err := statement.Query()
		if err != nil {
			return common.DataBaseError, err.Error(), nil
		}
		defer rows.Close()

		var resultSolutions []common.SolutionInfoItem
		for rows.Next() {
			var solutionInfo common.SolutionInfoItem
			rows.Scan(&solutionInfo.SolutionId, &solutionInfo.SolutionName)
			resultSolutions = append(resultSolutions, solutionInfo)
		}
		return common.Success, "", resultSolutions
	}
}

func (solutionMgr *SqliteSolutionManager) CreateNewSolution(solutionName string) (common.ErrorCode, string) {
	if !solutionMgr.isWriteable {
		panic("SqliteSolutionManager Need Writeable To CreateNewSolution")
	}

	querySQL := `SELECT count(*) FROM ai_solutions WHERE solutionName = ?`
	{
		statement, err := solutionMgr.sqliteDb.Prepare(querySQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		defer statement.Close()

		result := statement.QueryRow(solutionName)
		var count int
		result.Scan(&count)

		if count > 0 {
			return common.DuplicatedSolutionName, common.DuplicatedSolutionName.GetMsgFormat(solutionName)
		}
	}

	insertStudentSQL := `INSERT INTO ai_solutions(id, solutionName, solutionMeta) VALUES (?,?,?)`
	{
		statement, err := solutionMgr.sqliteDb.Prepare(insertStudentSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		defer statement.Close()

		_, err = statement.Exec(uuid.New().String(), solutionName, nil)
		if err != nil {
			return common.DataBaseError, err.Error()
		}
	}
	return common.Success, ""
}

//AssetSetManager

type SqliteAssetSetManager struct {
	sqliteDb    *sql.DB
	isWriteable bool
	locker      *sync.RWMutex
}

func (assetSetMgr *SqliteAssetSetManager) ListAssetSets(solutionId string) (common.ErrorCode, string, []common.AssetSetInfoItem) {
	querySQL := `SELECT id, solutionId, assetSetName FROM ai_asset_sets WHERE solutionId = ?`
	{
		statement, err := assetSetMgr.sqliteDb.Prepare(querySQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error(), nil
		}
		defer statement.Close()

		rows, err := statement.Query(solutionId)
		if err != nil {
			return common.DataBaseError, err.Error(), nil
		}
		defer rows.Close()

		var resultAssetSets []common.AssetSetInfoItem
		for rows.Next() {
			var assetSetInfo common.AssetSetInfoItem
			rows.Scan(&assetSetInfo.AssetSetId, &assetSetInfo.SolutionId, &assetSetInfo.AssetSetName)
			resultAssetSets = append(resultAssetSets, assetSetInfo)
		}
		return 0, "", resultAssetSets
	}
}

func (assetSetMgr *SqliteAssetSetManager) CreateAssetSet(solutionId string, assetSetName string) (common.ErrorCode, string) {
	if !assetSetMgr.isWriteable {
		panic("SqliteAssetSetManager Need Writeable To CreateAssetSet")
	}

	//Check Solution Exist
	{
		CheckSQL := `SELECT count(*) FROM ai_solutions WHERE id=?`
		statement, err := assetSetMgr.sqliteDb.Prepare(CheckSQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		defer statement.Close()

		result := statement.QueryRow(solutionId)
		var count int
		result.Scan(&count)
		if count == 0 {
			return common.InvalidSolution, common.InvalidSolution.GetMsgFormat(solutionId)
		}
	}

	//Check Duplicated AssetSet
	{
		CheckSQL := `SELECT count(*) FROM ai_asset_sets WHERE assetSetName = ?`
		statement, err := assetSetMgr.sqliteDb.Prepare(CheckSQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		defer statement.Close()

		result := statement.QueryRow(assetSetName)
		var count int
		result.Scan(&count)
		if count > 0 {
			return common.DuplicatedAssetSetName, common.DuplicatedAssetSetName.GetMsgFormat(assetSetName)
		}
	}

	//Do Insert
	{
		CreateSQL := `INSERT INTO ai_asset_sets(id, solutionId, assetSetName) VALUES (?,?,?)`
		statement, err := assetSetMgr.sqliteDb.Prepare(CreateSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			return common.DataBaseError, err.Error()
		}

		_, err = statement.Exec(uuid.New().String(), solutionId, assetSetName)
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		return common.Success, ""
	}
}

func (assetSetMgr *SqliteAssetSetManager) Release() {
	if assetSetMgr != nil {
		if assetSetMgr.isWriteable {
			assetSetMgr.locker.Unlock()
		} else {
			assetSetMgr.locker.RUnlock()
		}
	}
}

//AssetManager

type SqliteAssetManager struct {
	sqliteDb    *sql.DB
	isWriteable bool
	locker      *sync.RWMutex
}

func (assetMgr *SqliteAssetManager) ListAssets(assetSetId string) (common.ErrorCode, string, []common.AssetSummaryInfoItem) {
	querySQL := `SELECT id, assetSetId, assetType, assetName, assetVersion FROM ai_asset_documentations WHERE assetSetId = ?`
	statement, err := assetMgr.sqliteDb.Prepare(querySQL) // Prepare statement.
	if err != nil {
		return common.DataBaseError, err.Error(), nil
	}
	defer statement.Close()

	rows, err := statement.Query(assetSetId)
	if err != nil {
		return common.DataBaseError, err.Error(), nil
	}
	defer rows.Close()

	var resultAssets []common.AssetSummaryInfoItem
	for rows.Next() {
		var assetSetInfo common.AssetSummaryInfoItem
		rows.Scan(&assetSetInfo.AssetId, &assetSetInfo.AssetSetId, &assetSetInfo.AssetType, &assetSetInfo.AssetName, &assetSetInfo.AssetVersion)
		resultAssets = append(resultAssets, assetSetInfo)
	}
	return common.Success, "", resultAssets
}

func (assetMgr *SqliteAssetManager) CreateAsset(assetSetId string, assetType string, assetName string, assetInitContent string) (common.ErrorCode, string) {
	if !assetMgr.isWriteable {
		panic("SqliteAssetManager Need Writeable To CreateAsset")
	}
	//AssetSet Exist Check
	{
		querySQL := `SELECT count(*) FROM ai_asset_sets WHERE id = ?;`
		statement, err := assetMgr.sqliteDb.Prepare(querySQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		defer statement.Close()

		result := statement.QueryRow(assetSetId)
		var count int
		result.Scan(&count)

		if count != 1 {
			return common.InvalidAssetSet, common.InvalidAssetSet.GetMsgFormat(assetSetId)
		}
	}

	//AssetName Duplicated Check
	{
		querySQL := `SELECT count(*) FROM ai_asset_documentations WHERE assetName = ?;`
		statement, err := assetMgr.sqliteDb.Prepare(querySQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		defer statement.Close()

		result := statement.QueryRow(assetName)
		var count int
		result.Scan(&count)

		if count > 0 {
			return common.DuplicatedAssetName, common.DuplicatedAssetName.GetMsgFormat(assetName)
		}
	}

	{
		createSQL := `INSERT INTO ai_asset_documentations(id, assetSetId, assetType, assetName, assetContent, assetVersion) VALUES (?,?,?,?,?,?)`
		statement, err := assetMgr.sqliteDb.Prepare(createSQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		defer statement.Close()

		_, err = statement.Exec(uuid.New().String(), assetSetId, assetType, assetName, assetInitContent, uuid.New().String())
		if err != nil {
			return common.DataBaseError, err.Error()
		}
		return common.Success, ""
	}
}

func (assetMgr *SqliteAssetManager) Release() {
	if assetMgr != nil {
		if assetMgr.isWriteable {
			assetMgr.locker.Unlock()
		} else {
			assetMgr.locker.RUnlock()
		}
	}
}

//AssetDocument

type SqliteAssetDocument struct {
	assetId     string
	sqliteDb    *sql.DB
	isWriteable bool
	locker      *sync.RWMutex
}

func (dbDoc *SqliteAssetDocument) UpdateContent(content string) (common.ErrorCode, string, string) {
	if !dbDoc.isWriteable {
		panic("SqliteAssetDocument Need Writeable To UpdateContent")
	}

	newVersion := uuid.New().String()
	{
		updateSQL := `UPDATE ai_asset_documentations SET assetContent = ?, assetVersion = ? WHERE id = ?`

		statement, err := dbDoc.sqliteDb.Prepare(updateSQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error(), ""
		}

		_, err = statement.Exec(string(content), newVersion, dbDoc.assetId)
		if err != nil {
			return common.DataBaseError, err.Error(), ""
		}

		return common.Success, "", newVersion
	}
}

func (dbDoc *SqliteAssetDocument) ReadAsset() (common.ErrorCode, string, *common.AssetDetailInfo) {
	{
		querySQL := `SELECT id, assetSetId, assetType, assetName, assetContent, assetVersion FROM ai_asset_documentations WHERE id = ?`
		statement, err := dbDoc.sqliteDb.Prepare(querySQL) // Prepare statement.
		if err != nil {
			return common.DataBaseError, err.Error(), nil
		}

		result := statement.QueryRow(dbDoc.assetId)
		var content common.AssetDetailInfo
		result.Scan(&content.AssetId, &content.AssetSetId, &content.AssetType, &content.AssetName, &content.AssetContent, &content.AssetVersion)

		return common.Success, "", &content
	}
}

func (dbDoc *SqliteAssetDocument) Release() {
	if dbDoc != nil {
		if dbDoc.isWriteable {
			dbDoc.locker.Unlock()
		} else {
			dbDoc.locker.RUnlock()
		}
	}
}

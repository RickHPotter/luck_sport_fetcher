package models

const insertSQL = `
	INSERT OR IGNORE INTO Odds (MatchDate, League, HomeTeam, AwayTeam, EarlyOdds1, EarlyOddsX, EarlyOdds2, FinalOdds1, FinalOddsX, FinalOdds2, Score)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

const createTableOddsSQL = `
    CREATE TABLE IF NOT EXISTS Odds (
      ID INTEGER PRIMARY KEY AUTOINCREMENT,
      MatchDate TEXT NOT NULL,
      League TEXT,
      HomeTeam TEXT NOT NULL,
      AwayTeam TEXT NOT NULL,
      Score TEXT,
      EarlyOdds1 REAL,
      EarlyOddsX REAL,
      EarlyOdds2 REAL,
      FinalOdds1 REAL,
      FinalOddsX REAL,
      FinalOdds2 REAL,
      UNIQUE (HomeTeam, AwayTeam, MatchDate)
    );`

const createTableFutureOddsSQL = `
    CREATE TABLE IF NOT EXISTS FutureOdds (
      ID INTEGER PRIMARY KEY AUTOINCREMENT,
      MatchDate TEXT NOT NULL,
      League TEXT,
      HomeTeam TEXT NOT NULL,
      AwayTeam TEXT NOT NULL,
      Score TEXT,
      EarlyOdds1 REAL,
      EarlyOddsX REAL,
      EarlyOdds2 REAL,
      FinalOdds1 REAL,
      FinalOddsX REAL,
      FinalOdds2 REAL,
      UNIQUE (HomeTeam, AwayTeam, MatchDate)
    );`

package models

const insertSQL = `
	INSERT OR IGNORE INTO Odds (MatchDate, League, LeagueName, HomeTeam, AwayTeam, EarlyOdds1, EarlyOddsX, EarlyOdds2, FinalOdds1, FinalOddsX, FinalOdds2, Score)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

const schemaSQL = `
  CREATE TABLE IF NOT EXISTS Filter (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE,
    LeaguesChosen TEXT,
    YearsChosen TEXT,
    SameHomeAwayE INTEGER DEFAULT 0,
    SameHomeAwayF INTEGER DEFAULT 0,
    SameDrawE INTEGER DEFAULT 0,
    SameDrawF INTEGER DEFAULT 0,
    Early INTEGER DEFAULT 0,
    Final INTEGER DEFAULT 0,
    SameLeagueOnly INTEGER DEFAULT 0,
    OnlyTodayAndTomorrowMatches INTEGER DEFAULT 0,
    OnlyPresentEarlyOdds INTEGER DEFAULT 0,
    OnlyPresentFinalOdds INTEGER DEFAULT 0,
    HomeWinE1 TEXT,
    HomeWinE2 TEXT,
    HomeWinF1 TEXT,
    HomeWinF2 TEXT,
    DrawE1 TEXT,
    DrawE2 TEXT,
    DrawF1 TEXT,
    DrawF2 TEXT,
    AwayWinE1 TEXT,
    AwayWinE2 TEXT,
    AwayWinF1 TEXT,
    AwayWinF2 TEXT
  );

  CREATE TABLE IF NOT EXISTS League (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
  );

  CREATE TABLE IF NOT EXISTS Folder (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
  );

  CREATE TABLE IF NOT EXISTS LeagueFolder (
    FolderID INTEGER NOT NULL,
    LeagueID INTEGER NOT NULL,
    PRIMARY KEY (FolderID, LeagueID),
    FOREIGN KEY (FolderID) REFERENCES Folder(ID),
    FOREIGN KEY (LeagueID) REFERENCES League(ID)
  );

  CREATE TABLE IF NOT EXISTS Odds (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    MatchDate TEXT NOT NULL,
    League TEXT,
    LeagueID INTEGER,
    LeagueName TEXT,
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
    UNIQUE (HomeTeam, AwayTeam, MatchDate),
    FOREIGN KEY (LeagueID) REFERENCES League(ID)
  );

  CREATE TABLE IF NOT EXISTS FutureOdds (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    MatchDate TEXT NOT NULL,
    League TEXT,
    LeagueName TEXT,
    LeagueID INTEGER,
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
    FOREIGN KEY (LeagueID) REFERENCES League(ID)
  );
  `

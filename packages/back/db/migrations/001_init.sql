CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    role ENUM('superroot', 'admin', 'student') NOT NULL,
    assigned_sport VARCHAR(255)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tournaments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    sport VARCHAR(255) NOT NULL,
    weather_condition ENUM('sunny', 'rainy', 'any') NOT NULL
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE teams (
    id INT AUTO_INCREMENT PRIMARY KEY,
    class_id INT,
    name VARCHAR(255) NOT NULL,
    tournament_id INT,
    entry_status VARCHAR(255),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE matches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tournament_id INT,
    round INT,
    match_number_in_round INT,
    team1_id INT,
    team2_id INT,
    team1_score INT,
    team2_score INT,
    winner_team_id INT,
    next_match_id INT,
    status VARCHAR(255),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id),
    FOREIGN KEY (team1_id) REFERENCES teams(id),
    FOREIGN KEY (team2_id) REFERENCES teams(id),
    FOREIGN KEY (winner_team_id) REFERENCES teams(id),
    FOREIGN KEY (next_match_id) REFERENCES matches(id)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- チームごとのポイント内訳を記録するテーブル
CREATE TABLE team_points (
  id INT AUTO_INCREMENT PRIMARY KEY,
  class_id INT NOT NULL UNIQUE,
  init_score INT DEFAULT 0,
  attendance_score INT DEFAULT 0,
  volleyball1_score INT DEFAULT 0,
  volleyball2_score INT DEFAULT 0,
  volleyball3_score INT DEFAULT 0,
  volleyball_championship_score INT DEFAULT 0,
  table_tennis1_score INT DEFAULT 0,
  table_tennis2_score INT DEFAULT 0,
  table_tennis3_score INT DEFAULT 0,
  table_tennis_championship_score INT DEFAULT 0,
  table_tennis_rainy_bonus_score INT DEFAULT 0,
  soccer1_score INT DEFAULT 0,
  soccer2_score INT DEFAULT 0,
  soccer3_score INT DEFAULT 0,
  soccer_championship_score INT DEFAULT 0,
  relay_A_score INT DEFAULT 0,
  relay_B_score INT DEFAULT 0,
  relay_bonus_score INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
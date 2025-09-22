-- リレーシステムの作成（新しいclass_id構造対応）

-- 1. リレー結果テーブルの作成
CREATE TABLE relay_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    relay_type ENUM('A', 'B') NOT NULL COMMENT 'リレーブロック: A=Aブロック, B=Bブロック',
    relay_rank INT NOT NULL COMMENT '順位: 1-6位',
    class_id INT NOT NULL COMMENT 'クラスID: 0=未設定, 学年の代表クラスIDを使用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `idx_type_rank` (`relay_type`, `relay_rank`),
    CHECK (relay_rank BETWEEN 1 AND 6)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
COMMENT = 'リレー結果テーブル: 学年単位での順位管理';

-- 2. 初期データの挿入（未設定状態）
INSERT INTO relay_results (relay_type, relay_rank, class_id) VALUES
('A', 1, 0), ('A', 2, 0), ('A', 3, 0), ('A', 4, 0), ('A', 5, 0), ('A', 6, 0),
('B', 1, 0), ('B', 2, 0), ('B', 3, 0), ('B', 4, 0), ('B', 5, 0), ('B', 6, 0);

-- 3. team_pointsテーブルに新しいclass_id構造でデータを挿入
-- class_id構造: 1年生=11,12,13, 2年生=21,22,23, 3年生=31,32,33, 4年生=41,42,43, 5年生=51,52,53, 専・教=6
INSERT IGNORE INTO team_points (class_id, relay_A_score, relay_B_score, relay_bonus_score) VALUES
(11, 0, 0, 0), 
(12, 0, 0, 0),   
(13, 0, 0, 0), 
(21, 0, 0, 0),  
(23, 0, 0, 0),  
(31, 0, 0, 0),  
(32, 0, 0, 0),  
(33, 0, 0, 0),  
(41, 0, 0, 0),  
(42, 0, 0, 0),  
(43, 0, 0, 0),  
(51, 0, 0, 0),  
(52, 0, 0, 0),  
(53, 0, 0, 0),  
(6, 0, 0, 0);   
-- 4. class_nameカラムが存在しない場合は追加
SET @column_exists = (
    SELECT COUNT(*) 
    FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'team_points' 
    AND COLUMN_NAME = 'class_name'
);

SET @sql = IF(@column_exists = 0, 
    'ALTER TABLE team_points ADD COLUMN class_name VARCHAR(10) DEFAULT NULL AFTER class_id', 
    'SELECT "class_name column already exists" as message'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 5. クラス名データの更新
UPDATE team_points SET class_name = '1-1' WHERE class_id = 11;
UPDATE team_points SET class_name = '1-2' WHERE class_id = 12;
UPDATE team_points SET class_name = '1-3' WHERE class_id = 13;
UPDATE team_points SET class_name = 'IS2' WHERE class_id = 21;
UPDATE team_points SET class_name = 'IT2' WHERE class_id = 22;
UPDATE team_points SET class_name = 'IE2' WHERE class_id = 23;
UPDATE team_points SET class_name = 'IS3' WHERE class_id = 31;
UPDATE team_points SET class_name = 'IT3' WHERE class_id = 32;
UPDATE team_points SET class_name = 'IE3' WHERE class_id = 33;
UPDATE team_points SET class_name = 'IS4' WHERE class_id = 41;
UPDATE team_points SET class_name = 'IT4' WHERE class_id = 42;
UPDATE team_points SET class_name = 'IE4' WHERE class_id = 43;
UPDATE team_points SET class_name = 'IS5' WHERE class_id = 51;
UPDATE team_points SET class_name = 'IT5' WHERE class_id = 52;
UPDATE team_points SET class_name = 'IE5' WHERE class_id = 53;
UPDATE team_points SET class_name = '専・教' WHERE class_id = 6;

-- 6. 学年別得点集計ビューの作成
CREATE OR REPLACE VIEW relay_grade_scores AS
SELECT 
    CASE 
        WHEN class_id BETWEEN 11 AND 13 THEN 1
        WHEN class_id BETWEEN 21 AND 23 THEN 2
        WHEN class_id BETWEEN 31 AND 33 THEN 3
        WHEN class_id BETWEEN 41 AND 43 THEN 4
        WHEN class_id BETWEEN 51 AND 53 THEN 5
        WHEN class_id = 6 THEN 6
    END as grade,
    CASE 
        WHEN class_id BETWEEN 11 AND 13 THEN '1年生'
        WHEN class_id BETWEEN 21 AND 23 THEN '2年生'
        WHEN class_id BETWEEN 31 AND 33 THEN '3年生'
        WHEN class_id BETWEEN 41 AND 43 THEN '4年生'
        WHEN class_id BETWEEN 51 AND 53 THEN '5年生'
        WHEN class_id = 6 THEN '専・教'
    END as grade_name,
    MAX(relay_A_score) as relay_A_score,
    MAX(relay_B_score) as relay_B_score,
    MAX(relay_bonus_score) as relay_bonus_score,
    MAX(relay_A_score) + MAX(relay_B_score) as total_relay_score
FROM team_points 
WHERE class_id IN (11, 12, 13, 21, 22, 23, 31, 32, 33, 41, 42, 43, 51, 52, 53, 6)
GROUP BY grade
ORDER BY grade;

-- 7. インデックスの追加
-- idx_team_points_class_id
SET @exists := (
    SELECT COUNT(*) 
    FROM INFORMATION_SCHEMA.STATISTICS 
    WHERE table_schema = DATABASE()
      AND table_name = 'team_points'
      AND index_name = 'idx_team_points_class_id'
);
SET @sql = IF(@exists = 0,
    'CREATE INDEX idx_team_points_class_id ON team_points(class_id)',
    'SELECT "index idx_team_points_class_id already exists"'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- idx_team_points_class_name
SET @exists := (
    SELECT COUNT(*) 
    FROM INFORMATION_SCHEMA.STATISTICS 
    WHERE table_schema = DATABASE()
      AND table_name = 'team_points'
      AND index_name = 'idx_team_points_class_name'
);
SET @sql = IF(@exists = 0,
    'CREATE INDEX idx_team_points_class_name ON team_points(class_name)',
    'SELECT "index idx_team_points_class_name already exists"'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- idx_relay_results_type_rank
SET @exists := (
    SELECT COUNT(*) 
    FROM INFORMATION_SCHEMA.STATISTICS 
    WHERE table_schema = DATABASE()
      AND table_name = 'relay_results'
      AND index_name = 'idx_relay_results_type_rank'
);
SET @sql = IF(@exists = 0,
    'CREATE INDEX idx_relay_results_type_rank ON relay_results(relay_type, relay_rank)',
    'SELECT "index idx_relay_results_type_rank already exists"'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

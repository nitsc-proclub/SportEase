CREATE TABLE scoretable_settings (
    setting_key VARCHAR(255) PRIMARY KEY,
    setting_value TEXT NOT NULL
);

INSERT INTO scoretable_settings (setting_key, setting_value) VALUES ('showTotalScores', 'true');


CREATE TABLE weather_settings (
    setting_key VARCHAR(255) PRIMARY KEY,
    setting_value TEXT NOT NULL
);

INSERT INTO weather_settings (setting_key, setting_value) VALUES ('tableTennisWeather', 'sunny');
-- 卓球（雨天時）: 1回戦の敗者を敗者復活戦へ紐付ける
SET @tt_rainy_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時）');
SET @tt_rainy_loser1_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時・敗者戦左側）');
SET @tt_rainy_loser2_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時・敗者戦右側）');

-- 左側ブロックの敗者
SET @ttr_match_1_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 1);
SET @ttr_match_2_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 2);
SET @ttr_match_3_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 3);
SET @ttr_match_4_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 4);
SET @ttrl_match_13_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser1_tour_id AND match_number_in_round = 13);
UPDATE matches SET loser_next_match_id = @ttrl_match_13_id WHERE id IN (@ttr_match_1_id, @ttr_match_2_id, @ttr_match_3_id, @ttr_match_4_id);

-- 右側ブロックの敗者
SET @ttr_match_5_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 5);
SET @ttr_match_6_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 6);
SET @ttr_match_7_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 7);
SET @ttr_match_8_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_tour_id AND match_number_in_round = 8);
SET @ttrl_match_14_id = (SELECT id FROM matches WHERE tournament_id = @tt_rainy_loser2_tour_id AND match_number_in_round = 14);
UPDATE matches SET loser_next_match_id = @ttrl_match_14_id WHERE id IN (@ttr_match_5_id, @ttr_match_6_id, @ttr_match_7_id, @ttr_match_8_id);

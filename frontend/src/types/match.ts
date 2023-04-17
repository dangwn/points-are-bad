export type MatchWithoutGoals = {
  match_date: string,
  home: string,
  away: string
};

export type Match = {
  match_date: string,
  home: string,
  away: string,
  home_goals: number|null,
  away_goals: number|null
}
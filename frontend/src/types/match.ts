export type MatchWithoutGoals = {
  match_date: string,
  home: string,
  away: string
};

export type Match = MatchWithoutGoals & {
  home_goals: number|null,
  away_goals: number|null
}



export type MatchWithId = Match & {
  match_id: number
}
export type UserPrediction =   {
  prediction_id: number,
  home_goals: number|null,
  away_goals: number|null,
  match: {
    match_date: string,
    home: string,
    away: string,
    home_goals: number|null,
    away_goals: number|null
  },
  user: {
    username: string,
    is_admin: boolean
  }
};

export type NewPrediction = {
  prediction_id: number,
  home_goals: number|null,
  away_goals: number|null
};
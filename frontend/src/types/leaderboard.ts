export type LeaderboardUser = {
  user: {
    username: string
  },
  points: number,
  correct_scores: number,
  largest_error: number,
  position: number|null
}

export type LeaderboardApiResponse = {
  data: LeaderboardUser[],
  hasMore: boolean,
  pageIndex: number
}
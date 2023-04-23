export const CLIENT_HOST: string = process.env.CLIENT_HOST || 'http://localhost:3000';
export const API_HOST: string = process.env.API_HOST || 'http://localhost:8000';
export const NULL_PREDICTION_PENALTY: number = parseInt(process.env.NULL_PREDICTION_PENALTY || '10');

export const QUERY_OPTIONS = {
  refetchOnWindowFocus: false
}
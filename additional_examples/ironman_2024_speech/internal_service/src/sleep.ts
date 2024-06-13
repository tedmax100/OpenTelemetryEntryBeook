const sleep = async (timeInSeconds: number) => {
  return new Promise((resolve) => {
    setTimeout(resolve, timeInSeconds)
  })
}
export default sleep;
export interface Sleep {
  sleep: (ts: number) => Promise<unknown>
}
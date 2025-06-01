// Helper to imitate network delay
export const delay = (ms: number) => new Promise((res) => setTimeout(res, ms));

export const simulateUser = async (userId, apiList) => {
  const userResult = { userId, steps: [] };
  for (const { url, method } of apiList) {
    const start = performance.now();
    try {
      const res = await fetch(url, { method });
      const end = performance.now();
      userResult.steps.push({
        api: url,
        method,
        success: res.ok,
        status: res.status,
        time: end - start,
      });
    } catch (err) {
      const end = performance.now();
      userResult.steps.push({
        api: url,
        method,
        success: false,
        error: err.message,
        time: end - start,
      });
    }
  }
  return userResult;
};

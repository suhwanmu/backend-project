export const simulateUser = async (userId, apiList) => {
  const userResult = { userId, steps: [] };
  for (const { url, method } of apiList) {
    const start = performance.now();
    try {
      const res = await fetch(url, { method });
      const end = performance.now();

      const contentType = res.headers.get('content-type');
      let body;
      let success = res.ok;

      try {
        if (contentType?.includes('application/json')) {
          body = await res.json();
          // 👉 바디 내부에 실패 신호 있는지 검사 (예: "status": "fail", "error": ...)
          if (body?.status === 'fail' || body?.error) {
            success = false;
          }
        } else {
          body = await res.text(); // JSON이 아닐 경우 텍스트로
        }
      } catch (err) {
        console.warn('⚠️ 응답 바디 파싱 실패:', err.message);
        body = null;
        success = false;
      }

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

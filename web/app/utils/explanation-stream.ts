export type ExplanationEvent = { type: string; text?: string; result?: { answer: string; references?: unknown[]; coverage?: unknown } };

// JSON callers remain supported. A stream is successful only after a final
// validated result; a disconnect must never save provisional text as an answer.
export async function readExplanation(response: Response, onEvent: (event: ExplanationEvent) => void) {
  if (!response.headers?.get('Content-Type')?.includes('application/x-ndjson')) return response.json();
  if (!response.body) throw new Error('Explanation stream unavailable. Please retry.');
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';
  let result: ExplanationEvent['result'];
  function line(text: string) {
    if (!text.trim()) return;
    const event = JSON.parse(text) as ExplanationEvent;
    if (event.type === 'error') throw new Error(event.text || 'Explanation unavailable. Please retry.');
    if (event.type === 'result') {
      if (!event.result || typeof event.result.answer !== 'string' || !Array.isArray(event.result.references)) throw new Error('Invalid explanation response.');
      result = event.result;
    } else onEvent(event);
  }
  try {
    while (!result) {
      const { value, done } = await reader.read();
      buffer += decoder.decode(value, { stream: !done });
      if (buffer.length > 2 * 1024 * 1024) throw new Error('Explanation response is too large.');
      let end: number;
      while ((end = buffer.indexOf('\n')) >= 0) { line(buffer.slice(0, end)); buffer = buffer.slice(end + 1); }
      if (done) { line(buffer); break; }
    }
    if (!result) throw new Error('Explanation was interrupted. Please retry.');
    return result;
  } finally { await reader.cancel().catch(() => {}); reader.releaseLock(); }
}

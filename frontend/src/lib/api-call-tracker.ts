// API Call ID helper for frontend
// This helps track API calls from frontend side

/**
 * Get the API call ID from a response
 * Returns the X-API-Call-ID header value
 *
 * @param response - Fetch API Response object
 * @returns The API call ID or null if not present
 *
 * @example
 * const response = await fetch('/api/v1/billing/plans');
 * const callID = getAPICallID(response);
 * console.log(`API Call: ${callID}`);
 */
export function getAPICallID(response: Response): string | null {
  return response.headers.get('X-API-Call-ID');
}

/**
 * Enhanced fetch wrapper that logs API call IDs in development
 *
 * @param url - Request URL
 * @param options - Fetch options
 * @returns Fetch response
 *
 * @example
 * // In development, automatically logs call IDs
 * const response = await fetchWithCallID('/api/v1/billing/plans');
 */
export async function fetchWithCallID(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  const response = await fetch(url, options);

  // Log call ID in development
  if (import.meta.env?.DEV || window?.location?.hostname === 'localhost') {
    const callID = getAPICallID(response);
    if (callID) {
      console.log(`🔍 API Call [${callID.substring(0, 8)}] ${options.method || 'GET'} ${url} → ${response.status}`);
    }
  }

  return response;
}

/**
 * Track API calls for debugging
 * Stores recent calls in memory
 */
class APICallTracker {
  private calls: Array<{
    callID: string;
    method: string;
    url: string;
    status: number;
    timestamp: Date;
  }> = [];

  private maxCalls = 50; // Keep last 50 calls

  /**
   * Track a response
   */
  track(response: Response, method: string = 'GET', url: string = '') {
    const callID = getAPICallID(response);
    if (callID) {
      this.calls.push({
        callID,
        method,
        url,
        status: response.status,
        timestamp: new Date(),
      });

      // Keep only recent calls
      if (this.calls.length > this.maxCalls) {
        this.calls.shift();
      }
    }
  }

  /**
   * Get recent API calls
   */
  getCalls() {
    return [...this.calls];
  }

  /**
   * Find API call by ID
   */
  findCall(callID: string) {
    return this.calls.find(call => call.callID === callID);
  }

  /**
   * Clear tracking history
   */
  clear() {
    this.calls = [];
  }

  /**
   * Export calls as JSON
   */
  exportJSON() {
    return JSON.stringify(this.calls, null, 2);
  }
}

// Export singleton instance
export const apiTracker = new APICallTracker();

/**
 * Console command helper for debugging
 * Use in browser console:
 *
 * window.trackAPI() // Start tracking
 * window.showAPICalls() // Show recent calls
 * window.exportAPICalls() // Export to clipboard
 */
if (typeof window !== 'undefined') {
  (window as any).trackAPI = () => {
    console.log('🔍 API Call tracking enabled!');
    console.log('Use window.showAPICalls() to see recent calls');
    console.log('Use window.exportAPICalls() to export to console');
  };

  (window as any).showAPICalls = () => {
    const calls = apiTracker.getCalls();
    console.table(calls.map(call => ({
      ID: call.callID.substring(0, 8),
      Method: call.method,
      URL: call.url,
      Status: call.status,
      Time: call.timestamp.toLocaleTimeString(),
    })));
    return calls;
  };

  (window as any).exportAPICalls = () => {
    const json = apiTracker.exportJSON();
    console.log(json);
    return json;
  };
}

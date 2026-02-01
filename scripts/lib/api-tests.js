const { buildIfNeeded, logError, logInfo, startServer, stopProcess, waitForHealth } = require('./runner');

const testCategories = {
  health: 'Health check endpoints',
  list: 'List items endpoints',
  crud: 'CRUD operations',
  validation: 'Validation error tests',
};

function generateTests(selected) {
  const tests = [];
  const errorSchema = {
    type: 'object',
    required: ['error'],
    properties: {
      error: {
        type: 'object',
        required: ['code', 'message', 'details'],
      },
    },
  };
  const itemImageSchema = {
    type: 'object',
    required: ['url'],
    properties: {
      url: { type: 'string' },
      alt: { type: 'string' },
      sort_order: { type: 'integer', minimum: 0 },
    },
  };
  const itemSchema = {
    type: 'object',
    required: [
      'id',
      'title',
      'price_cents',
      'condition',
      'status',
      'is_featured',
      'country',
      'delivery_available',
      'created_at',
      'updated_at',
      'images',
    ],
    properties: {
      id: { type: 'integer', minimum: 1 },
      title: { type: 'string', minLength: 3 },
      description: { type: 'string', nullable: true },
      price_cents: { type: 'integer', minimum: 0 },
      category: { type: 'string', nullable: true },
      condition: { type: 'string', enum: ['new', 'like_new', 'good', 'fair', 'parts', 'unknown'] },
      status: { type: 'string', enum: ['draft', 'active', 'reserved', 'sold', 'archived'] },
      is_featured: { type: 'boolean' },
      city: { type: 'string', nullable: true },
      postal_code: { type: 'string', nullable: true },
      country: { type: 'string' },
      delivery_available: { type: 'boolean' },
      created_at: { type: 'string', format: 'date-time' },
      updated_at: { type: 'string', format: 'date-time' },
      published_at: { type: 'string', format: 'date-time', nullable: true },
      images: { type: 'array', items: itemImageSchema },
    },
  };
  const listItemsResponseSchema = {
    type: 'object',
    required: ['items'],
    properties: {
      items: { type: 'array', items: itemSchema },
      next_cursor: { type: 'string', format: 'byte', nullable: true },
    },
  };

  if (selected.includes('health')) {
    tests.push({
      category: 'health',
      name: 'GET /health returns 200',
      method: 'GET',
      path: '/health',
      expectedStatus: 200,
      expectedSchema: {
        type: 'object',
        required: ['status', 'timestamp'],
      },
    });
  }

  if (selected.includes('list')) {
    tests.push({
      category: 'list',
      name: 'GET /v1/items returns list response',
      method: 'GET',
      path: '/v1/items',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with status filter',
      method: 'GET',
      path: '/v1/items?status=active',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with category filter',
      method: 'GET',
      path: '/v1/items?category=electronics',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with price range filter',
      method: 'GET',
      path: '/v1/items?min_price_cents=1000&max_price_cents=50000',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with location filter',
      method: 'GET',
      path: '/v1/items?city=Strasbourg&postal_code=67000',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with featured/delivery filters',
      method: 'GET',
      path: '/v1/items?is_featured=true&delivery_available=false',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with sort',
      method: 'GET',
      path: '/v1/items?sort=created_at:desc',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with limit',
      method: 'GET',
      path: '/v1/items?limit=10',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with cursor',
      method: 'GET',
      path: '/v1/items?cursor=eyJpZCI6MjB9',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with full-text search',
      method: 'GET',
      path: '/v1/items?q=iPhone',
      expectedStatus: 200,
      expectedSchema: listItemsResponseSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with invalid status returns 400',
      method: 'GET',
      path: '/v1/items?status=invalid_status',
      expectedStatus: 400,
      expectedSchema: errorSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with limit below range returns 400',
      method: 'GET',
      path: '/v1/items?limit=0',
      expectedStatus: 400,
      expectedSchema: errorSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with limit above range returns 400',
      method: 'GET',
      path: '/v1/items?limit=101',
      expectedStatus: 400,
      expectedSchema: errorSchema,
    });

    tests.push({
      category: 'list',
      name: 'GET /v1/items with invalid sort returns 400',
      method: 'GET',
      path: '/v1/items?sort=created_at:down',
      expectedStatus: 400,
      expectedSchema: errorSchema,
    });
  }

  if (selected.includes('crud')) {
    tests.push({
      category: 'crud',
      name: 'GET /v1/items/999 returns 404',
      method: 'GET',
      path: '/v1/items/999',
      expectedStatus: 404,
      expectedSchema: {
        type: 'object',
        required: ['error'],
      },
    });
  }

  if (selected.includes('validation')) {
    tests.push({
      category: 'validation',
      name: 'POST /v1/items with short title returns 400',
      method: 'POST',
      path: '/v1/items',
      body: {
        title: 'AB',
        price_cents: 10000,
        condition: 'good',
      },
      expectedStatus: 400,
      expectedSchema: {
        type: 'object',
        required: ['error'],
      },
    });

    tests.push({
      category: 'validation',
      name: 'POST /v1/items with negative price returns 400',
      method: 'POST',
      path: '/v1/items',
      body: {
        title: 'Test Item',
        price_cents: -100,
        condition: 'good',
      },
      expectedStatus: 400,
      expectedSchema: {
        type: 'object',
        required: ['error'],
      },
    });

    tests.push({
      category: 'validation',
      name: 'POST /v1/items with invalid condition returns 400',
      method: 'POST',
      path: '/v1/items',
      body: {
        title: 'Test Item',
        price_cents: 10000,
        condition: 'invalid_condition',
      },
      expectedStatus: 400,
      expectedSchema: {
        type: 'object',
        required: ['error'],
      },
    });
  }

  return tests;
}

function validateSchema(data, schema, path = '') {
  const errors = [];
  const currentPath = path || 'response';

  if (!schema || typeof schema !== 'object') {
    return errors;
  }

  if (schema.nullable && data === null) {
    return errors;
  }

  if (schema.oneOf) {
    const matches = schema.oneOf.some(candidate => validateSchema(data, candidate, currentPath).length === 0);
    if (!matches) {
      errors.push(`${currentPath}: value does not match any allowed schema`);
    }
    return errors;
  }

  if (schema.type) {
    if (schema.type === 'array') {
      if (!Array.isArray(data)) {
        errors.push(`${currentPath}: expected array but got ${typeof data}`);
        return errors;
      }
    } else if (schema.type === 'object') {
      if (!data || typeof data !== 'object' || Array.isArray(data)) {
        errors.push(`${currentPath}: expected object but got ${Array.isArray(data) ? 'array' : typeof data}`);
        return errors;
      }
    } else if (schema.type === 'integer') {
      if (!Number.isInteger(data)) {
        errors.push(`${currentPath}: expected integer but got ${typeof data}`);
        return errors;
      }
    } else if (typeof data !== schema.type) {
      errors.push(`${currentPath}: expected ${schema.type} but got ${typeof data}`);
      return errors;
    }
  }

  if (schema.enum && !schema.enum.includes(data)) {
    errors.push(`${currentPath}: expected one of ${JSON.stringify(schema.enum)} but got ${JSON.stringify(data)}`);
  }

  if (typeof data === 'string') {
    if (schema.minLength !== undefined && data.length < schema.minLength) {
      errors.push(`${currentPath}: expected length >= ${schema.minLength} but got ${data.length}`);
    }
    if (schema.format === 'date-time') {
      const isoRegex = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z$/;
      if (!isoRegex.test(data)) {
        errors.push(`${currentPath}: expected ISO 8601 date-time but got ${JSON.stringify(data)}`);
      }
    }
  }

  if (typeof data === 'number') {
    if (schema.minimum !== undefined && data < schema.minimum) {
      errors.push(`${currentPath}: expected >= ${schema.minimum} but got ${data}`);
    }
    if (schema.maximum !== undefined && data > schema.maximum) {
      errors.push(`${currentPath}: expected <= ${schema.maximum} but got ${data}`);
    }
  }

  if (schema.type === 'array' && schema.items && Array.isArray(data)) {
    for (let i = 0; i < data.length; i++) {
      errors.push(...validateSchema(data[i], schema.items, `${currentPath}[${i}]`));
    }
  }

  if ((schema.type === 'object' || schema.properties) && data && typeof data === 'object' && !Array.isArray(data)) {
    if (schema.required) {
      for (const key of schema.required) {
        if (!(key in data)) {
          errors.push(`${currentPath}: missing required key '${key}'`);
        }
      }
    }

    if (schema.properties) {
      for (const [key, propertySchema] of Object.entries(schema.properties)) {
        if (key in data) {
          errors.push(...validateSchema(data[key], propertySchema, `${currentPath}.${key}`));
        }
      }
    }
  }

  return errors;
}

async function makeRequest(method, url, headers = {}, body = null) {
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
  };

  if (body) {
    options.body = JSON.stringify(body);
  }

  const response = await fetch(url, options);

  let data;
  const contentType = response.headers.get('content-type');
  if (contentType && contentType.includes('application/json')) {
    data = await response.json();
  } else {
    data = await response.text();
  }

  return {
    status: response.status,
    headers: Object.fromEntries(response.headers.entries()),
    data,
  };
}

async function runTest(server, test, options) {
  const url = `http://localhost:${server.port}${test.path}`;

  if (options.verbose) {
    logInfo(`    ${test.method} ${test.path}`);
  }

  try {
    const response = await makeRequest(test.method, url, {}, test.body || null);

    if (response.status !== test.expectedStatus) {
      return {
        passed: false,
        reason: `status code mismatch: expected ${test.expectedStatus} but got ${response.status}`,
        expected: { status: test.expectedStatus },
        actual: { status: response.status, data: response.data },
      };
    }

    if (test.expectedSchema) {
      const schemaErrors = validateSchema(response.data, test.expectedSchema);
      if (schemaErrors.length > 0) {
        return {
          passed: false,
          reason: `schema validation failed: ${schemaErrors.join(', ')}`,
          expected: test.expectedSchema,
          actual: response.data,
        };
      }
    }

    return { passed: true, response };
  } catch (error) {
    return {
      passed: false,
      reason: `network error: ${error.message}`,
      error: error.message,
    };
  }
}

function printTestResult(server, test, result, options, state) {
  if (options.tap) {
    if (result.passed) {
      console.log(`ok ${state.tapTestNumber++} - ${server.id}: ${test.name}`);
    } else {
      console.log(`not ok ${state.tapTestNumber++} - ${server.id}: ${test.name}`);
      console.log(`  # ${result.reason}`);
      if (options.verbose && result.actual) {
        console.log(`  # Expected: ${JSON.stringify(result.expected)}`);
        console.log(`  # Actual: ${JSON.stringify(result.actual)}`);
      }
    }
    return;
  }

  if (result.passed) {
    logInfo(`  ✓ ${test.name}`);
  } else {
    logError(`  ✗ ${test.name}`);
    logError(`    Reason: ${result.reason}`);
    if (options.verbose && result.actual) {
      logInfo(`    Expected: ${JSON.stringify(result.expected)}`);
      logInfo(`    Actual: ${JSON.stringify(result.actual)}`);
    }
  }
}

async function runTestsForServer(server, tests, options, state) {
  logInfo(`\n${server.label} (${server.port})`);
  logInfo('='.repeat(50));

  const serverResults = {
    total: 0,
    passed: 0,
    failed: 0,
    tests: [],
  };

  await buildIfNeeded(server);
  const child = await startServer(server);
  const ready = await waitForHealth(server, options.timeoutMs);

  if (!ready) {
    await stopProcess(child, server);
    for (const test of tests) {
      const result = { passed: false, reason: 'Server health check failed' };
      printTestResult(server, test, result, options, state);
      serverResults.total++;
      serverResults.failed++;
      serverResults.tests.push({ test, result });
      state.results.total++;
      state.results.failed++;
    }
    state.results.byServer[server.id] = serverResults;
    return serverResults;
  }

  for (const test of tests) {
    const result = await runTest(server, test, options);
    printTestResult(server, test, result, options, state);

    serverResults.total++;
    if (result.passed) {
      serverResults.passed++;
      state.results.passed++;
    } else {
      serverResults.failed++;
      state.results.failed++;
      state.results.failures.push({
        server: server.id,
        test: test.name,
        ...result,
      });
    }
    serverResults.tests.push({ test, result });
    state.results.total++;
  }

  await stopProcess(child, server);
  state.results.byServer[server.id] = serverResults;
  return serverResults;
}

function printSummary(state, options) {
  if (options.tap) {
    console.log(`1..${state.results.total}`);
    console.log(`# ${state.results.passed} passed, ${state.results.failed} failed`);
    return;
  }

  console.log('');
  console.log('='.repeat(60));
  logInfo('Summary');
  console.log('='.repeat(60));

  logInfo(`Total tests: ${state.results.total}`);
  logInfo(`Passed: ${state.results.passed}`);
  logInfo(`Failed: ${state.results.failed}`);

  console.log('');
  logInfo('Results by server:');
  for (const [serverId, results] of Object.entries(state.results.byServer)) {
    const status = results.failed === 0 ? '✓' : '✗';
    logInfo(`  ${status} ${serverId}: ${results.passed}/${results.total} passed`);
  }

  if (state.results.failures.length > 0) {
    console.log('');
    logError('Failed tests:');
    for (const failure of state.results.failures) {
      logError(`  - ${failure.server}: ${failure.test}`);
      logError(`    ${failure.reason}`);
    }
  }
}

async function writeJsonReport(state, output) {
  if (!output) {
    return;
  }

  const report = {
    timestamp: new Date().toISOString(),
    summary: {
      total: state.results.total,
      passed: state.results.passed,
      failed: state.results.failed,
    },
    byServer: state.results.byServer,
    failures: state.results.failures,
  };

  const fs = await import('fs/promises');
  await fs.writeFile(output, JSON.stringify(report, null, 2));
  logInfo(`\nJSON report written to ${output}`);
}

async function runApiVerification(servers, options) {
  const tests = generateTests(options.tests);
  const state = {
    tapTestNumber: 1,
    results: {
      total: 0,
      passed: 0,
      failed: 0,
      byServer: {},
      failures: [],
    },
  };

  if (options.tap) {
    console.log('TAP version 13');
  }

  logInfo(`\nGenerated ${tests.length} tests`);

  for (const server of servers) {
    await runTestsForServer(server, tests, options, state);
  }

  printSummary(state, options);
  await writeJsonReport(state, options.output);

  if (state.results.failed > 0) {
    process.exit(1);
  }
}

module.exports = {
  runApiVerification,
  testCategories,
};

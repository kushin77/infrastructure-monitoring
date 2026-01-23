const http = require('http');

const options = {
    hostname: 'localhost',
    port: 8081,
    path: '/health',
    method: 'GET',
    timeout: 5000
};

const req = http.request(options, (res) => {
    if (res.statusCode === 200) {
        console.log('OK');
        process.exit(0);
    } else {
        console.log('NOT OK');
        process.exit(1);
    }
});

req.on('error', (err) => {
    console.log('NOT OK');
    process.exit(1);
});

req.on('timeout', () => {
    console.log('NOT OK');
    process.exit(1);
});

req.end();

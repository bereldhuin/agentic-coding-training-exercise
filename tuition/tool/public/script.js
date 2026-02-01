document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('register-form');
    const serverList = document.getElementById('server-list');
    const getIpBtn = document.getElementById('get-ip-btn');
    const ipOutput = document.getElementById('ip-output');

    // Fetch and display servers
    const fetchServers = async () => {
        try {
            const response = await fetch('/api/servers');
            const servers = await response.json();
            renderServers(servers);
        } catch (error) {
            console.error('Error fetching servers:', error);
        }
    };

    const renderServers = (servers) => {
        serverList.innerHTML = servers.map(server => `
            <li class="${server.status}">
                <div class="server-info">
                    <strong>${server.name}</strong><br>
                    <small>${server.url}</small><br>
                    <span class="status-badge status-${server.status}">${server.status}</span>
                    ${server.lastChecked ? `<small>Last checked: ${new Date(server.lastChecked).toLocaleString()}</small>` : ''}
                </div>
                <div class="server-actions">
                    <button onclick="verifyServer('${server.id}')">Verify</button>
                </div>
            </li>
        `).join('');
    };

    // Register new server
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('name').value;
        const url = document.getElementById('url').value;

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, url })
            });
            if (response.ok) {
                registerForm.reset();
                fetchServers();
            }
        } catch (error) {
            console.error('Error registering server:', error);
        }
    });

    // Verify server
    window.verifyServer = async (id) => {
        try {
            const response = await fetch(`/api/verify/${id}`, { method: 'POST' });
            if (response.ok) {
                fetchServers();
            }
        } catch (error) {
            console.error('Error verifying server:', error);
        }
    };

    // Get IP Address
    getIpBtn.addEventListener('click', async () => {
        ipOutput.textContent = 'Running command...';
        try {
            const response = await fetch('/api/ip-cmd');
            const data = await response.json();
            if (data.output) {
                ipOutput.textContent = data.output;
            } else if (data.error) {
                ipOutput.textContent = 'Error: ' + data.error;
            } else {
                ipOutput.textContent = 'No output';
            }
        } catch (error) {
            ipOutput.textContent = 'Failed to execute command';
        }
    });

    // Initial load
    fetchServers();
});

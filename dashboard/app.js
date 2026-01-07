const API_BASE = '/api';
let token = localStorage.getItem('moderator_token');
let userRole = localStorage.getItem('moderator_role') || 'full';

// Helper functions
async function apiCall(endpoint, method = 'GET', body = null) {
    const headers = {
        'Content-Type': 'application/json',
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const options = { method, headers };
    if (body) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(`${API_BASE}${endpoint}`, options);
    const data = await response.json();

    if (response.status === 401) {
        logout();
        throw new Error('Session expired');
    }

    return data;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('it-IT', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Auth functions
async function login(e) {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorEl = document.getElementById('login-error');

    try {
        const data = await apiCall('/login', 'POST', { username, password });

        if (data.status === 'ok') {
            token = data.token;
            userRole = data.role || 'full';
            localStorage.setItem('moderator_token', token);
            localStorage.setItem('moderator_role', userRole);
            showDashboard();
        } else {
            errorEl.textContent = data.msg || 'Credenziali non valide';
        }
    } catch (err) {
        errorEl.textContent = 'Errore di connessione';
    }
}

function logout() {
    token = null;
    userRole = 'full';
    localStorage.removeItem('moderator_token');
    localStorage.removeItem('moderator_role');
    document.getElementById('login-page').classList.remove('hidden');
    document.getElementById('dashboard-page').classList.add('hidden');

    // Reset menu visibility
    document.querySelectorAll('.nav-links li').forEach(li => li.classList.remove('hidden'));
}

async function checkAuth() {
    if (!token) return false;
    try {
        const data = await apiCall('/verify');
        return data.status === 'ok';
    } catch {
        return false;
    }
}

async function showDashboard() {
    document.getElementById('login-page').classList.add('hidden');
    document.getElementById('dashboard-page').classList.remove('hidden');

    // Check if user has limited access (users_only)
    if (userRole === 'users_only') {
        // Hide all menu items except "Utenti"
        document.querySelectorAll('.nav-links li').forEach(li => {
            const link = li.querySelector('a');
            if (link && link.dataset.section !== 'users') {
                li.classList.add('hidden');
            }
        });

        // Show only users section
        document.querySelectorAll('.section').forEach(s => s.classList.add('hidden'));
        document.getElementById('users-section').classList.remove('hidden');

        // Set users link as active
        document.querySelectorAll('.nav-links a').forEach(a => a.classList.remove('active'));
        const usersLink = document.querySelector('.nav-links a[data-section="users"]');
        if (usersLink) usersLink.classList.add('active');
    } else {
        // Full access - load cities as default
        await loadCities();
    }
}

// Navigation
function setupNavigation() {
    document.querySelectorAll('.nav-links a').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const section = e.target.dataset.section;

            // Update active link
            document.querySelectorAll('.nav-links a').forEach(l => l.classList.remove('active'));
            e.target.classList.add('active');

            // Show section
            document.querySelectorAll('.section').forEach(s => s.classList.add('hidden'));
            document.getElementById(`${section}-section`).classList.remove('hidden');

            // Load data
            switch(section) {
                case 'cities': loadCities(); break;
                case 'schools': loadSchools(); break;
                case 'posts': loadPendingPosts(); loadReportedPosts(); loadAllPosts(); break;
                case 'spotted': loadPendingSpotted(); loadReportedSpotted(); loadAllSpotted(); break;
                case 'users': break; // Users loaded on search
            }
        });
    });

    // Tabs
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const tab = e.target.dataset.tab;
            const parent = e.target.closest('.section');

            parent.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
            e.target.classList.add('active');

            parent.querySelectorAll('.tab-content').forEach(c => c.classList.add('hidden'));
            parent.querySelector(`#${tab}`).classList.remove('hidden');
        });
    });
}

// Cities
async function loadCities() {
    try {
        const data = await apiCall('/cities');
        const tbody = document.querySelector('#cities-table tbody');
        tbody.innerHTML = '';

        if (!data.data || data.data.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="empty-state">Nessuna città trovata</td></tr>';
            return;
        }

        data.data.forEach(city => {
            tbody.innerHTML += `
                <tr>
                    <td>${city.id}</td>
                    <td>${city.name}</td>
                    <td>${city.region || '-'}</td>
                    <td class="actions">
                        <button class="btn btn-secondary btn-small" onclick="editCity(${city.id}, '${city.name}', '${city.region || ''}')">Modifica</button>
                        <button class="btn btn-danger btn-small" onclick="confirmDeleteCity(${city.id})">Elimina</button>
                    </td>
                </tr>
            `;
        });

        // Update city selects
        updateCitySelects(data.data);
    } catch (err) {
        console.error('Error loading cities:', err);
    }
}

function updateCitySelects(cities) {
    const schoolCityFilter = document.getElementById('school-city-filter');
    const schoolCity = document.getElementById('school-city');

    const optionsHtml = cities.map(c => `<option value="${c.id}">${c.name}</option>`).join('');

    schoolCityFilter.innerHTML = '<option value="">Tutte le città</option>' + optionsHtml;
    schoolCity.innerHTML = '<option value="">Seleziona città</option>' + optionsHtml;
}

function showAddCityModal() {
    document.getElementById('city-modal-title').textContent = 'Aggiungi Città';
    document.getElementById('city-id').value = '';
    document.getElementById('city-name').value = '';
    document.getElementById('city-region').value = '';
    document.getElementById('city-modal').classList.remove('hidden');
}

function editCity(id, name, region) {
    document.getElementById('city-modal-title').textContent = 'Modifica Città';
    document.getElementById('city-id').value = id;
    document.getElementById('city-name').value = name;
    document.getElementById('city-region').value = region;
    document.getElementById('city-modal').classList.remove('hidden');
}

function closeCityModal() {
    document.getElementById('city-modal').classList.add('hidden');
}

async function saveCity(e) {
    e.preventDefault();
    const id = document.getElementById('city-id').value;
    const name = document.getElementById('city-name').value;
    const region = document.getElementById('city-region').value;

    try {
        if (id) {
            await apiCall(`/cities/${id}`, 'PUT', { name, region });
        } else {
            await apiCall('/cities', 'POST', { name, region });
        }
        closeCityModal();
        await loadCities();
    } catch (err) {
        alert('Errore nel salvataggio');
    }
}

function confirmDeleteCity(id) {
    document.getElementById('confirm-message').textContent = 'Sei sicuro di voler eliminare questa città? Tutte le scuole associate verranno rimosse.';
    document.getElementById('confirm-delete-btn').onclick = () => deleteCity(id);
    document.getElementById('confirm-modal').classList.remove('hidden');
}

async function deleteCity(id) {
    try {
        await apiCall(`/cities/${id}`, 'DELETE');
        closeConfirmModal();
        await loadCities();
    } catch (err) {
        alert('Errore nell\'eliminazione. Assicurati che non ci siano scuole associate.');
    }
}

// Schools
async function loadSchools() {
    try {
        const cityId = document.getElementById('school-city-filter').value;
        const endpoint = cityId ? `/schools?city_id=${cityId}` : '/schools';
        const data = await apiCall(endpoint);
        const tbody = document.querySelector('#schools-table tbody');
        tbody.innerHTML = '';

        if (!data.data || data.data.length === 0) {
            tbody.innerHTML = '<tr><td colspan="5" class="empty-state">Nessuna scuola trovata</td></tr>';
            return;
        }

        data.data.forEach(school => {
            tbody.innerHTML += `
                <tr>
                    <td>${school.id}</td>
                    <td>${school.name}</td>
                    <td>${school.city_name}</td>
                    <td>${school.email_domain || '-'}</td>
                    <td class="actions">
                        <button class="btn btn-secondary btn-small" onclick="editSchool(${school.id}, '${school.name}', '${school.email_domain || ''}')">Modifica</button>
                        <button class="btn btn-danger btn-small" onclick="confirmDeleteSchool(${school.id})">Elimina</button>
                    </td>
                </tr>
            `;
        });
    } catch (err) {
        console.error('Error loading schools:', err);
    }
}

function showAddSchoolModal() {
    document.getElementById('school-modal-title').textContent = 'Aggiungi Scuola';
    document.getElementById('school-id').value = '';
    document.getElementById('school-name').value = '';
    document.getElementById('school-city').value = '';
    document.getElementById('school-city').disabled = false;
    document.getElementById('school-domain').value = '';
    document.getElementById('school-modal').classList.remove('hidden');
}

function editSchool(id, name, domain) {
    document.getElementById('school-modal-title').textContent = 'Modifica Scuola';
    document.getElementById('school-id').value = id;
    document.getElementById('school-name').value = name;
    document.getElementById('school-city').disabled = true;
    document.getElementById('school-domain').value = domain;
    document.getElementById('school-modal').classList.remove('hidden');
}

function closeSchoolModal() {
    document.getElementById('school-modal').classList.add('hidden');
}

async function saveSchool(e) {
    e.preventDefault();
    const id = document.getElementById('school-id').value;
    const name = document.getElementById('school-name').value;
    const cityId = parseInt(document.getElementById('school-city').value);
    const emailDomain = document.getElementById('school-domain').value;

    try {
        if (id) {
            await apiCall(`/schools/${id}`, 'PUT', { name, email_domain: emailDomain });
        } else {
            await apiCall('/schools', 'POST', { name, city_id: cityId, email_domain: emailDomain });
        }
        closeSchoolModal();
        await loadSchools();
    } catch (err) {
        alert('Errore nel salvataggio');
    }
}

function confirmDeleteSchool(id) {
    document.getElementById('confirm-message').textContent = 'Sei sicuro di voler eliminare questa scuola?';
    document.getElementById('confirm-delete-btn').onclick = () => deleteSchool(id);
    document.getElementById('confirm-modal').classList.remove('hidden');
}

async function deleteSchool(id) {
    try {
        await apiCall(`/schools/${id}`, 'DELETE');
        closeConfirmModal();
        await loadSchools();
    } catch (err) {
        alert('Errore nell\'eliminazione. Assicurati che non ci siano utenti associati.');
    }
}

// Posts
async function loadPendingPosts() {
    try {
        const data = await apiCall('/posts/pending');
        const container = document.getElementById('pending-posts-list');

        if (!data.data || data.data.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>Nessun post in attesa di approvazione</p></div>';
            return;
        }

        container.innerHTML = data.data.map(post => `
            <div class="card">
                <div class="card-header">
                    <div>
                        <strong>${post.creator_first_name} ${post.creator_last_name}</strong>
                        <span class="badge badge-warning">In attesa</span>
                    </div>
                </div>
                <div class="card-meta">
                    <span>${post.school_name || 'N/A'} - ${post.city_name || 'N/A'}</span>
                    <span>${formatDate(post.creation_timestamp)}</span>
                </div>
                <div class="card-content">${post.content}</div>
                <div class="card-actions">
                    <button class="btn btn-success btn-small" onclick="approvePost(${post.id})">Approva</button>
                    <button class="btn btn-danger btn-small" onclick="rejectPost(${post.id})">Rifiuta</button>
                    <button class="btn btn-secondary btn-small" onclick="confirmDeletePost(${post.id})">Elimina</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        console.error('Error loading pending posts:', err);
    }
}

async function loadReportedPosts() {
    try {
        const data = await apiCall('/posts/reported');
        const container = document.getElementById('reported-posts-list');

        if (!data.data || data.data.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>Nessun post segnalato</p></div>';
            return;
        }

        container.innerHTML = data.data.map(post => `
            <div class="card">
                <div class="card-header">
                    <div>
                        <strong>${post.creator_first_name} ${post.creator_last_name}</strong>
                        <span class="badge badge-danger">${post.report_count} segnalazioni</span>
                    </div>
                </div>
                <div class="card-meta">
                    <span>${post.school_name || 'N/A'} - ${post.city_name || 'N/A'}</span>
                    <span>${formatDate(post.creation_timestamp)}</span>
                </div>
                <div class="card-content">${post.content}</div>
                <div class="card-actions">
                    <button class="btn btn-danger btn-small" onclick="confirmDeletePost(${post.id})">Elimina</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        console.error('Error loading reported posts:', err);
    }
}

async function approvePost(id) {
    try {
        await apiCall(`/posts/${id}/approve`, 'PUT');
        await loadPendingPosts();
    } catch (err) {
        alert('Errore nell\'approvazione');
    }
}

async function rejectPost(id) {
    try {
        await apiCall(`/posts/${id}/reject`, 'PUT');
        await loadPendingPosts();
    } catch (err) {
        alert('Errore nel rifiuto');
    }
}

function confirmDeletePost(id) {
    document.getElementById('confirm-message').textContent = 'Sei sicuro di voler eliminare questo post?';
    document.getElementById('confirm-delete-btn').onclick = () => deletePost(id);
    document.getElementById('confirm-modal').classList.remove('hidden');
}

async function deletePost(id) {
    try {
        await apiCall(`/posts/${id}`, 'DELETE');
        closeConfirmModal();
        await loadPendingPosts();
        await loadReportedPosts();
        await loadAllPosts();
    } catch (err) {
        alert('Errore nell\'eliminazione');
    }
}

function getStatusBadge(status) {
    const badges = {
        'received': '<span class="badge badge-warning">In attesa</span>',
        'approved': '<span class="badge badge-success">Approvato</span>',
        'rejected': '<span class="badge badge-danger">Rifiutato</span>'
    };
    return badges[status] || status;
}

async function loadAllPosts() {
    try {
        const data = await apiCall('/posts');
        const container = document.getElementById('all-posts-list');

        if (!data.data || data.data.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>Nessun post trovato</p></div>';
            return;
        }

        container.innerHTML = data.data.map(post => `
            <div class="card">
                <div class="card-header">
                    <div>
                        <strong>${post.creator_first_name} ${post.creator_last_name}</strong>
                        ${getStatusBadge(post.status)}
                    </div>
                </div>
                <div class="card-meta">
                    <span>Email: ${post.creator_email}</span>
                    <span>${post.school_name || 'N/A'} - ${post.city_name || 'N/A'}</span>
                    <span>${formatDate(post.creation_timestamp)}</span>
                </div>
                <div class="card-content">${post.content}</div>
                <div class="card-actions">
                    <select onchange="setPostStatus(${post.id}, this.value)" class="status-select">
                        <option value="">Cambia stato...</option>
                        <option value="received" ${post.status === 'received' ? 'disabled' : ''}>In attesa</option>
                        <option value="approved" ${post.status === 'approved' ? 'disabled' : ''}>Approvato</option>
                        <option value="rejected" ${post.status === 'rejected' ? 'disabled' : ''}>Rifiutato</option>
                    </select>
                    <button class="btn btn-danger btn-small" onclick="confirmDeletePost(${post.id})">Elimina</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        console.error('Error loading all posts:', err);
    }
}

async function setPostStatus(id, status) {
    if (!status) return;
    try {
        await apiCall(`/posts/${id}/status`, 'PUT', { status });
        await loadAllPosts();
        await loadPendingPosts();
    } catch (err) {
        alert('Errore nel cambio stato');
    }
}

// Spotted
async function loadPendingSpotted() {
    try {
        const data = await apiCall('/spotted/pending');
        const container = document.getElementById('pending-spotted-list');

        if (!data.data || data.data.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>Nessuno spotted in attesa di approvazione</p></div>';
            return;
        }

        container.innerHTML = data.data.map(s => `
            <div class="card" style="border-left: 4px solid ${s.color || '#6366f1'}">
                <div class="card-header">
                    <div>
                        <strong>${s.creator_first_name} ${s.creator_last_name}</strong>
                        <span class="badge badge-warning">${s.visibility_desc}</span>
                    </div>
                </div>
                <div class="card-meta">
                    <span>Email: ${s.creator_email}</span>
                    <span>${s.school_name || 'N/A'} - ${s.city_name || 'N/A'}</span>
                    <span>${formatDate(s.creation_timestamp)}</span>
                </div>
                <div class="card-content">${s.content}</div>
                <div class="card-actions">
                    <button class="btn btn-success btn-small" onclick="approveSpotted(${s.id})">Approva</button>
                    <button class="btn btn-danger btn-small" onclick="rejectSpotted(${s.id})">Rifiuta</button>
                    <button class="btn btn-secondary btn-small" onclick="confirmDeleteSpotted(${s.id})">Elimina</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        console.error('Error loading pending spotted:', err);
    }
}

async function loadReportedSpotted() {
    try {
        const data = await apiCall('/spotted/reported');
        const container = document.getElementById('reported-spotted-list');

        if (!data.data || data.data.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>Nessuno spotted segnalato</p></div>';
            return;
        }

        container.innerHTML = data.data.map(s => `
            <div class="card" style="border-left: 4px solid ${s.color || '#6366f1'}">
                <div class="card-header">
                    <div>
                        <strong>${s.creator_first_name} ${s.creator_last_name}</strong>
                        <span class="badge badge-danger">${s.report_count} segnalazioni</span>
                    </div>
                </div>
                <div class="card-meta">
                    <span>Email: ${s.creator_email}</span>
                    <span>${s.school_name || 'N/A'} - ${s.city_name || 'N/A'}</span>
                    <span>${formatDate(s.creation_timestamp)}</span>
                </div>
                <div class="card-content">${s.content}</div>
                <div class="card-actions">
                    <button class="btn btn-danger btn-small" onclick="confirmDeleteSpotted(${s.id})">Elimina</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        console.error('Error loading reported spotted:', err);
    }
}

async function approveSpotted(id) {
    try {
        await apiCall(`/spotted/${id}/approve`, 'PUT');
        await loadPendingSpotted();
    } catch (err) {
        alert('Errore nell\'approvazione');
    }
}

async function rejectSpotted(id) {
    try {
        await apiCall(`/spotted/${id}/reject`, 'PUT');
        await loadPendingSpotted();
    } catch (err) {
        alert('Errore nel rifiuto');
    }
}

function confirmDeleteSpotted(id) {
    document.getElementById('confirm-message').textContent = 'Sei sicuro di voler eliminare questo spotted?';
    document.getElementById('confirm-delete-btn').onclick = () => deleteSpotted(id);
    document.getElementById('confirm-modal').classList.remove('hidden');
}

async function deleteSpotted(id) {
    try {
        await apiCall(`/spotted/${id}`, 'DELETE');
        closeConfirmModal();
        await loadPendingSpotted();
        await loadReportedSpotted();
        await loadAllSpotted();
    } catch (err) {
        alert('Errore nell\'eliminazione');
    }
}

async function loadAllSpotted() {
    try {
        const data = await apiCall('/spotted');
        const container = document.getElementById('all-spotted-list');

        if (!data.data || data.data.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>Nessuno spotted trovato</p></div>';
            return;
        }

        container.innerHTML = data.data.map(s => `
            <div class="card" style="border-left: 4px solid ${s.color || '#6366f1'}">
                <div class="card-header">
                    <div>
                        <strong>${s.creator_first_name} ${s.creator_last_name}</strong>
                        ${getStatusBadge(s.status)}
                        <span class="badge">${s.visibility_desc}</span>
                    </div>
                </div>
                <div class="card-meta">
                    <span>Email: ${s.creator_email}</span>
                    <span>${s.school_name || 'N/A'} - ${s.city_name || 'N/A'}</span>
                    <span>${formatDate(s.creation_timestamp)}</span>
                </div>
                <div class="card-content">${s.content}</div>
                <div class="card-actions">
                    <select onchange="setSpottedStatus(${s.id}, this.value)" class="status-select">
                        <option value="">Cambia stato...</option>
                        <option value="received" ${s.status === 'received' ? 'disabled' : ''}>In attesa</option>
                        <option value="approved" ${s.status === 'approved' ? 'disabled' : ''}>Approvato</option>
                        <option value="rejected" ${s.status === 'rejected' ? 'disabled' : ''}>Rifiutato</option>
                    </select>
                    <button class="btn btn-danger btn-small" onclick="confirmDeleteSpotted(${s.id})">Elimina</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        console.error('Error loading all spotted:', err);
    }
}

async function setSpottedStatus(id, status) {
    if (!status) return;
    try {
        await apiCall(`/spotted/${id}/status`, 'PUT', { status });
        await loadAllSpotted();
        await loadPendingSpotted();
    } catch (err) {
        alert('Errore nel cambio stato');
    }
}

function closeConfirmModal() {
    document.getElementById('confirm-modal').classList.add('hidden');
}

// Users
function handleUserSearchKeyup(event) {
    if (event.key === 'Enter') {
        searchUsers();
    }
}

async function searchUsers() {
    const searchTerm = document.getElementById('user-search-input').value.trim();
    const resultsDiv = document.getElementById('users-results');
    const table = document.getElementById('users-table');
    const tbody = document.querySelector('#users-table tbody');

    if (searchTerm.length < 2) {
        resultsDiv.innerHTML = '<div class="empty-state"><p>Inserisci almeno 2 caratteri per cercare</p></div>';
        resultsDiv.classList.remove('hidden');
        table.classList.add('hidden');
        return;
    }

    try {
        const data = await apiCall(`/users/search?q=${encodeURIComponent(searchTerm)}`);

        if (!data.data || data.data.length === 0) {
            resultsDiv.innerHTML = '<div class="empty-state"><p>Nessun utente trovato</p></div>';
            resultsDiv.classList.remove('hidden');
            table.classList.add('hidden');
            return;
        }

        tbody.innerHTML = '';
        data.data.forEach(user => {
            const isRepresentative = user.role === 'representative';
            const roleBadgeClass = isRepresentative ? 'badge-role-representative' : 'badge-role-user';
            const roleLabel = isRepresentative ? 'Rappresentante' : 'Utente';
            const buttonText = isRepresentative ? 'Rimuovi Rappresentante' : 'Rendi Rappresentante';
            const buttonClass = isRepresentative ? 'btn-danger' : 'btn-primary';

            tbody.innerHTML += `
                <tr>
                    <td>${user.id}</td>
                    <td>${user.first_name} ${user.last_name}</td>
                    <td>
                        <div>${user.email}</div>
                        ${user.personal_email ? `<div style="color:#888;font-size:0.85rem">${user.personal_email}</div>` : ''}
                    </td>
                    <td>${user.school_name || 'N/A'}${user.city_name ? ` (${user.city_name})` : ''}</td>
                    <td><span class="badge ${roleBadgeClass}">${roleLabel}</span></td>
                    <td class="actions">
                        <button class="btn ${buttonClass} btn-small" onclick="toggleUserRole(${user.id}, '${user.role}')">${buttonText}</button>
                    </td>
                </tr>
            `;
        });

        resultsDiv.classList.add('hidden');
        table.classList.remove('hidden');
    } catch (err) {
        console.error('Error searching users:', err);
        resultsDiv.innerHTML = '<div class="empty-state"><p>Errore nella ricerca</p></div>';
        resultsDiv.classList.remove('hidden');
        table.classList.add('hidden');
    }
}

async function toggleUserRole(userId, currentRole) {
    const newRole = currentRole === 'representative' ? 'user' : 'representative';
    const action = newRole === 'representative' ? 'rendere Rappresentante' : 'rimuovere il ruolo Rappresentante a';

    if (!confirm(`Sei sicuro di voler ${action} questo utente?`)) {
        return;
    }

    try {
        await apiCall(`/users/${userId}/role`, 'PUT', { role: newRole });
        await searchUsers(); // Refresh results
    } catch (err) {
        alert('Errore nel cambio ruolo');
    }
}

// Initialize
document.addEventListener('DOMContentLoaded', async () => {
    document.getElementById('login-form').addEventListener('submit', login);
    document.getElementById('logout-btn').addEventListener('click', logout);
    document.getElementById('city-form').addEventListener('submit', saveCity);
    document.getElementById('school-form').addEventListener('submit', saveSchool);

    setupNavigation();

    if (await checkAuth()) {
        showDashboard();
    }
});

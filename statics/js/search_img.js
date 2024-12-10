function isMobile() {
    return window.innerWidth <= 768;
}

// 全局变量，用于跟踪当前标签
let currentTagName = null;
let pageState = {
    currentPage: 1,
    pageSize: 10,
    total: 0,
    searchType: 'all', // 'all', 'name', 'tag'
    order: 'desc'  // 默认降序
};

// 添加分页控件渲染函数
function renderPagination() {
    const totalPages = Math.max(1, Math.ceil(pageState.total / pageState.pageSize));
    const paginationHTML = `
        <div class="pagination">
            <button ${pageState.currentPage === 1 ? 'disabled' : ''} 
                    onclick="changePage(${pageState.currentPage - 1})">上一页</button>
            <span>第 ${pageState.currentPage} 页，共 ${totalPages} 页 (共${pageState.total}条记录)</span>
            <button ${pageState.currentPage === totalPages || totalPages === 0 ? 'disabled' : ''} 
                    onclick="changePage(${pageState.currentPage + 1})">下一页</button>
        </div>
    `;
    const paginationContainer = document.createElement('div');
    paginationContainer.innerHTML = paginationHTML;

    // 确保分页控件始终显示
    const resultArea = document.getElementById('resultArea');
    resultArea.appendChild(paginationContainer);
}

// 修改 toggleTimeSort 函数
function toggleTimeSort() {
    pageState.order = pageState.order === 'asc' ? 'desc' : 'asc';

    if (currentTagName) {
        // 如果当前是标签搜索状态，只对标签搜索结果进行排序
        searchByTag(currentTagName);
    } else {
        // 否则是普通搜索
        const isAllImages = document.getElementById('imageName').value === '';
        fetchImages(isAllImages);
    }
}

// 修改 searchByTag 函数
async function searchByTag(tagName) {
    console.log('搜索标签:', tagName); // 调试日志
    currentTagName = tagName;
    pageState.currentPage = 1;
    pageState.searchType = 'tag';

    const token = localStorage.getItem('token');
    if (!token) {
        console.error('No token found');
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/searchbytag', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                tags: [tagName],
                page: pageState.currentPage,
                pageSize: pageState.pageSize,
                order: pageState.order
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const responseData = await response.json();
        console.log('标签搜索响应:', responseData); // 调试日志
        displayResults(responseData);

        // 高亮当前选中的标签
        document.querySelectorAll('.tag').forEach(tag => {
            tag.classList.remove('active');
            if (tag.textContent.trim() === tagName) {
                tag.classList.add('active');
            }
        });
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('resultArea').textContent = '搜索过程中发生错误，请稍后再试。';
    }
}

// 将函数声明改为异步
async function fetchImages(isAllImages) {
    currentTagName = null;
    pageState.searchType = 'all';

    const searchInput = document.getElementById('imageName').value;

    if (!isAllImages && searchInput.length === 0) {
        alert('搜索内容不能为空！');
        return;
    }

    const token = localStorage.getItem('token');
    const url = isAllImages ? 'http://localhost:8080/searchAllimg' : 'http://localhost:8080/searchimg';
    const payload = {
        ...(isAllImages ? { allimg: true } : { name: searchInput }),
        page: pageState.currentPage,
        pageSize: pageState.pageSize,
        order: pageState.order
    };

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const responseData = await response.json();
        displayResults(responseData);
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('resultArea').textContent = '搜索过程中发生错误，请稍后再试。';
    }
}

// 修改 changePage 函数
function changePage(newPage) {
    pageState.currentPage = newPage;
    if (currentTagName) {
        // 如果是标签搜索状态，调用标签搜索
        searchByTag(currentTagName);
    } else {
        // 否则是普通搜索
        const isAllImages = document.getElementById('imageName').value === '';
        fetchImages(isAllImages);
    }
}

// 修改 displayResults 函数
function displayResults(responseData) {
    const resultArea = document.getElementById('resultArea');
    resultArea.innerHTML = '';

    pageState.total = responseData.data.total || 0;
    let data = responseData.data.images || responseData.data;

    if (data && Array.isArray(data) && data.length > 0) {
        if (isMobile()) {
            const tableHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>图片名称</th>
                            <th class="sortable" onclick="toggleTimeSort()">上传时间 ${pageState.order === 'asc' ? '↑' : '↓'}</th>
                            <th>标签</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${data.map(image => {
                console.log('处理的图片数据:', image);
                return `
                            <tr>
                                <td><a href="${image.image_url || ''}" target="_blank">${image.image_name || ''}</a></td>
                                <td>${image.UploadTime ? new Date(image.UploadTime).toLocaleString() : ''}</td>
                                <td>${Array.isArray(image.Tags) && image.Tags.length > 0 ?
                        image.Tags.map(tag =>
                            `<span class="tag" onclick="searchByTag('${tag.TagName}')">${tag.TagName}</span>`
                        ).join(' ')
                        : '无标签'}</td>
                            </tr>
                        `}).join('')}
                    </tbody>
                </table>
            `;
            resultArea.innerHTML = tableHTML;
        } else {
            const tableHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>图片ID</th>
                            <th>图片名称</th>
                            <th>哈希</th>
                            <th>大小</th>
                            <th>类型</th>
                            <th class="sortable" onclick="toggleTimeSort()">上传时间 ${pageState.order === 'asc' ? '↑' : '↓'}</th>
                            <th>描述</th>
                            <th>标签</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${data.map(image => {
                console.log('处理的图片数据:', image);
                return `
                            <tr>
                                <td>${image.id || ''}</td>
                                <td><a href="${image.image_url || ''}" target="_blank">${image.image_name || ''}</a></td>
                                <td>${image.hash_image || ''}</td>
                                <td>${image.image_size || ''}</td>
                                <td>${image.image_type || ''}</td>
                                <td>${image.UploadTime ? new Date(image.UploadTime).toLocaleString() : ''}</td>
                                <td>${image.description || ''}</td>
                                <td>${Array.isArray(image.Tags) && image.Tags.length > 0 ?
                        image.Tags.map(tag =>
                            `<span class="tag" onclick="searchByTag('${tag.TagName}')">${tag.TagName}</span>`
                        ).join(' ')
                        : '无标签'}</td>
                            </tr>
                        `}).join('')}
                    </tbody>
                </table>
            `;
            resultArea.innerHTML = tableHTML;
        }
    } else {
        resultArea.textContent = '没有找到匹配的图片。';
    }

    // 始终显示分页控件
    renderPagination();
}

// 修改事件监听部分
document.addEventListener('DOMContentLoaded', function () {
    // 搜索按钮点击事件
    document.getElementById('searchButton').addEventListener('click', function () {
        pageState.currentPage = 1;  // 重置分页
        fetchImages(false);
    });

    // 获取所有图片按钮点击事件
    document.getElementById('allImagesButton').addEventListener('click', function () {
        pageState.currentPage = 1;  // 重置分页
        fetchImages(true);
    });

    // 标签点击事件
    document.addEventListener('click', async function (e) {
        if (e.target.classList.contains('tag')) {
            const tagName = e.target.textContent;
            const token = localStorage.getItem('token');

            // 重置分页状态
            pageState.currentPage = 1;

            try {
                const response = await fetch('http://localhost:8080/searchbytag', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify({
                        tags: [tagName],
                        page: pageState.currentPage,
                        pageSize: pageState.pageSize
                    })
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const responseData = await response.json();
                console.log('Tag search response:', responseData); // 调试日志

                // 直接使用原始响应数据
                displayResults(responseData);
            } catch (error) {
                console.error('Error:', error);
                document.getElementById('resultArea').textContent = '搜索过程中发生错误，请稍后再试。';
            }
        }
    });

    // 添加标签搜索按钮事件
    document.getElementById('tagSearchButton').addEventListener('click', getAllTags);

    // 初始加载标签
    getAllTags();

    // 移除重复的标签点击事件监听器
    // 因为 searchByTag 函数已经在标签的 onclick 中直接调用
});

// 添加获取所有标签的函数
async function getAllTags() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch('http://localhost:8080/getalltag', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        displayTags(data.data);
    } catch (error) {
        console.error('Error fetching tags:', error);
        document.getElementById('tagsContainer').innerHTML = '获取标签失败';
    }
}

// 修改标签显示函数
function displayTags(tags) {
    const tagsContainer = document.getElementById('tagsContainer');
    if (!tags || tags.length === 0) {
        tagsContainer.innerHTML = '<div class="no-tags">暂无标签</div>';
        return;
    }

    // 按标签名称排序
    tags.sort((a, b) => a.TagName.localeCompare(b.TagName));

    tagsContainer.innerHTML = tags.map(tag => `
        <span class="tag" data-tagname="${tag.TagName}">
            ${tag.TagName}
        </span>
    `).join('');

    // 使用事件委托添加点击事件
    tagsContainer.addEventListener('click', function (e) {
        if (e.target.classList.contains('tag')) {
            e.preventDefault(); // 阻止默认行为
            e.stopPropagation(); // 阻止事件冒泡
            const tagName = e.target.dataset.tagname;
            if (tagName) {
                searchByTag(tagName);
            }
        }
    });
}

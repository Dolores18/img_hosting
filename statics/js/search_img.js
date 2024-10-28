function isMobile() {
    return window.innerWidth <= 768;
}

// 添加分页状态管理
const pageState = {
    currentPage: 1,
    pageSize: 10,
    total: 0
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

// 修改 displayResults 函数
function displayResults(responseData) {
    const resultArea = document.getElementById('resultArea');
    resultArea.innerHTML = ''; // 清空现有内容

    // 更新分页状态
    pageState.total = responseData.data.total || 0;  // 如果没有 total，默认为 0

    let data = responseData.data.images || responseData.data;

    if (data && Array.isArray(data) && data.length > 0) {
        if (isMobile()) {
            const tableHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>图片名称</th>
                            <th>上传时间</th>
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
                            <th>上传时间</th>
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

// 修改 fetchImages 函数
async function fetchImages(isAllImages) {
    const searchInput = document.getElementById('imageName').value;

    // 添加搜索长度检查
    if (!isAllImages && searchInput.length === 0) {
        alert('搜索内容不能为空！');
        return;
    }

    const token = localStorage.getItem('token');
    const url = isAllImages ? 'http://localhost:8080/searchAllimg' : 'http://localhost:8080/searchimg';
    const payload = {
        ...(isAllImages ? { allimg: true } : { name: searchInput }),
        page: pageState.currentPage,
        pageSize: pageState.pageSize
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

// 添加页码切换函数
function changePage(newPage) {
    pageState.currentPage = newPage;
    // 重新获取当前查询类型的数据
    const isAllImages = document.getElementById('imageName').value === '';
    fetchImages(isAllImages);
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
});

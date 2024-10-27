function isMobile() {
    return window.innerWidth <= 768;
}

function displayResults(responseData) {
    const resultArea = document.getElementById('resultArea');
    console.log('Response data:', responseData); // 调试日志

    // 处理不同的数据结构
    let data;
    if (responseData.data.images) {
        // 处理带有 total 的格式
        data = responseData.data.images;
    } else {
        // 处理直接返回数组的格式
        data = responseData.data;
    }

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
}

async function fetchImages(isAllImages) {
    const token = localStorage.getItem('token');
    const url = isAllImages ? 'http://localhost:8080/searchAllimg' : 'http://localhost:8080/searchimg';
    const payload = isAllImages ? { allimg: true } : { name: document.getElementById('imageName').value };

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

document.getElementById('searchForm').addEventListener('submit', function (e) {
    e.preventDefault();
    fetchImages(false);
});

document.getElementById('allImagesButton').addEventListener('click', function () {
    fetchImages(true);
});

// 修改标签点击事件处理
document.addEventListener('click', async function (e) {
    if (e.target.classList.contains('tag')) {
        const tagName = e.target.textContent;
        const token = localStorage.getItem('token');

        try {
            // 修改为正确的 API 端点
            const response = await fetch('http://localhost:8080/searchbytag', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ tags: [tagName] })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const responseData = await response.json();
            console.log('Tag search response:', responseData); // 调试日志

            // 处理嵌套的数据结构
            const modifiedData = {
                data: responseData.data.images || []
            };

            displayResults(modifiedData);
        } catch (error) {
            console.error('Error:', error);
            document.getElementById('resultArea').textContent = '搜索过程中发生错误，请稍后再试。';
        }
    }
});

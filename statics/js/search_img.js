function displayResults(responseData) {
    const resultArea = document.getElementById('resultArea');
    if (responseData.data && Array.isArray(responseData.data) && responseData.data.length > 0) {
        const tableHTML = `
        <table>
            <thead>
                <tr>
                    <th>图片ID</th>
                    <th>用户ID</th>
                    <th>图片URL</th>
                    <th>图片名称</th>
                    <th>扩展名</th>
                    <th>哈希</th>
                    <th>大小</th>
                    <th>类型</th>
                    <th>上传时间</th>
                    <th>描述</th>
                </tr>
            </thead>
            <tbody>
                ${responseData.data.map(image => `
                    <tr>
                        <td title="${image.ImageID || ''}">${image.ImageID || ''}</td>
                        <td title="${image.UserID || ''}">${image.UserID || ''}</td>
                        <td title="${image.ImageURL || ''}"><a href="${image.ImageURL || ''}" target="_blank">${image.ImageURL || ''}</a></td>
                        <td title="${image.ImageName || ''}">${image.ImageName || ''}</td>
                        <td title="${image.Imageextenion || ''}">${image.Imageextenion || ''}</td>
                        <td title="${image.HashImage || ''}">${image.HashImage || ''}</td>
                        <td title="${image.ImageSize || ''}">${image.ImageSize || ''}</td>
                        <td title="${image.ImageType || ''}">${image.ImageType || ''}</td>
                        <td title="${image.UploadTime ? new Date(image.UploadTime).toLocaleString() : ''}">${image.UploadTime ? new Date(image.UploadTime).toLocaleString() : ''}</td>
                        <td title="${image.Description || ''}">${image.Description || ''}</td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
        `;
        resultArea.innerHTML = tableHTML;
    } else {
        resultArea.textContent = '没有找到匹配的图片。';
    }
}

async function fetchImages(isAllImages) {
    const token = localStorage.getItem('token');
    const url = isAllImages ? 'https://imgapi.3049589.xyz/searchAllimg' : 'https://imgapi.3049589.xyz/searchimg';
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
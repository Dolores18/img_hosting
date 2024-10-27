// 存储所有要上传的文件
let filesArray = [];

// 监听粘贴事件
document.addEventListener('paste', function (event) {
    event.preventDefault();
    const items = event.clipboardData.items;

    for (let i = 0; i < items.length; i++) {
        const item = items[i];
        // 检查是否是图片
        if (item.type.indexOf('image') !== -1) {
            const file = item.getAsFile();
            addFileToArray(file);
        }
    }
});

function addFileToArray(file) {
    filesArray.push(file);
    previewFile(file);
}

function previewFile(file) {
    const preview = document.getElementById('image-preview');
    const reader = new FileReader();

    reader.onload = e => {
        const img = document.createElement('img');
        img.src = e.target.result;
        img.classList.add('thumb');
        preview.appendChild(img);
    };

    reader.readAsDataURL(file);
}

function previewImages() {
    const preview = document.getElementById('image-preview');
    const files = document.getElementById('file-input').files;

    if (files) {
        [...files].forEach(file => {
            addFileToArray(file);
        });
    }
}

async function uploadFiles() {
    const token = localStorage.getItem('token');
    const tags = document.getElementById('tag-input').value
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag.length > 0);

    if (filesArray.length > 0 && token) {
        const formData = new FormData();
        filesArray.forEach(file => {
            formData.append('upload[]', file);
        });

        try {
            const response = await fetch('http://localhost:8080/imgupload', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            if (response.ok) {
                const json = await response.json();
                console.log('Files uploaded successfully');

                // 如果有标签，为图片添加标签
                if (tags.length > 0) {
                    await addTagsToImage(json.data.imageid, tags);
                }

                alert(json.msg);
                clearFiles();
            } else {
                const json = await response.json();
                console.error('Failed to upload files');
                alert(json.error);
            }
        } catch (error) {
            console.error('Error:', error);
            alert('上传失败: ' + error.message);
        }
    } else {
        console.error('No files selected or token missing');
        alert('未选择文件或缺少token');
    }
}

async function addTagsToImage(imageId, tags) {
    const token = localStorage.getItem('token');

    for (const tagname of tags) {
        try {
            const response = await fetch('http://localhost:8080/addimagetag', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    imageid: imageId,
                    tagname: tagname
                })
            });

            if (!response.ok) {
                console.error(`Failed to add tag: ${tagname}`);
            }
        } catch (error) {
            console.error(`Error adding tag ${tagname}:`, error);
        }
    }
}

function clearFiles() {
    document.getElementById('file-input').value = '';
    document.getElementById('image-preview').innerHTML = '';
    document.getElementById('tag-input').value = ''; // 清空标签输入
    filesArray = [];
}

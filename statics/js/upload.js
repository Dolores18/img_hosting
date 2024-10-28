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
    // 获取原始输入并打印
    const rawInput = document.getElementById('tag-input').value;
    console.log('原始标签输入:', rawInput);  // 添加日志

    // 处理中英文空格，并分割标签
    const tags = rawInput
        .replace(/\u3000/g, ' ')  // 将中文空格替换为英文空格
        .split(/\s+/)             // 用一个或多个空格分割
        .map(tag => tag.trim())
        .filter(tag => tag.length > 0);

    console.log('处理后的标签数组:', tags);  // 添加日志
    alert('将添加以下标签：\n' + tags.join('\n')); // 弹窗显示要添加的标签

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

                alert('上传成功！\n添加的标签：' + tags.join(' ')); // 成功后显示添加的标签
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

// 修改标签添加函数
async function addTagsToImage(imageId, tags) {
    const token = localStorage.getItem('token');

    try {
        const response = await fetch('http://localhost:8080/addimagetag', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                imageid: imageId,
                tagnames: tags
            })
        });

        // 修改这部分错误处理逻辑
        const data = await response.json();

        if (!response.ok) {
            // 只有在真正的错误时才提示
            console.error('Failed to add tags:', data.error);
            alert('添加标签失败: ' + data.error);
            return false;
        }

        return true;  // 成功添加标签

    } catch (error) {
        console.error('Error adding tags:', error);
        alert('添加标签失败: ' + error.message);
        return false;
    }
}

function clearFiles() {
    document.getElementById('file-input').value = '';
    document.getElementById('image-preview').innerHTML = '';
    document.getElementById('tag-input').value = ''; // 清空标签输入
    filesArray = [];
}

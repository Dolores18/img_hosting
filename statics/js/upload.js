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

    if (filesArray.length > 0 && token) {
        const formData = new FormData();
        filesArray.forEach(file => {
            formData.append('upload[]', file);
        });

        try {
            const response = await fetch('https://imgapi.3049589.xyz/imgupload', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            if (response.ok) {
                const json = await response.json();
                console.log('Files uploaded successfully');
                alert(json.msg, json.data);
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

function clearFiles() {
    document.getElementById('file-input').value = '';
    document.getElementById('image-preview').innerHTML = '';
    filesArray = []; // 清空文件数组
}
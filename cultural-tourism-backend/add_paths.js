const fs = require('fs');
const path = require('path');

// ================= 配置区域 =================
const ROOT_DIR = process.cwd();

// 支持添加注释的文件后缀及对应的注释符号
const EXT_MAP = {
    '.go':   '//',
    '.js':   '//',
    '.ts':   '//',
    '.java': '//',
    '.c':    '//',
    '.cpp':  '//',
    '.php':  '//',
    '.sh':   '#',
    '.py':   '#',
    '.yml':  '#',
    '.yaml': '#',
    '.env':  '#',
    '.gitignore': '#'
    // 注意：JSON 不支持标准注释，Markdown 建议手动维护，故不在此列
};

// 忽略目录/文件 (正则)
const IGNORE_PATTERNS = [
    /^\./,              // 忽略 .git, .idea
    /^node_modules$/,
    /^tmp$/,
    /^docs$/,           // swagger docs 也是生成的，最好别动
    /^vendor$/,
    /^tests$/,
    /\.exe$/,
    /^go\.sum$/,
    /^go\.mod$/,        // go.mod 加注释可能会影响某些工具解析
    /^add_paths\.js$/,  // 忽略自己
    /^scan_project\.js$/
];
// ===========================================

function shouldIgnore(name) {
    return IGNORE_PATTERNS.some(pattern => pattern.test(name));
}

// 规范化路径分隔符 (把 Windows 的 \ 换成 /)
function normalizePath(p) {
    return p.split(path.sep).join('/');
}

function processFile(filePath) {
    const ext = path.extname(filePath);
    const commentSymbol = EXT_MAP[ext];

    // 如果是不支持的文件类型，直接跳过
    if (!commentSymbol) return;

    try {
        const content = fs.readFileSync(filePath, 'utf8');
        const lines = content.split('\n');
        
        // 计算相对路径，例如: controllers/region_controller.go
        const relativePath = normalizePath(path.relative(ROOT_DIR, filePath));
        
        // 构造标准注释行: // File: controllers/region_controller.go
        const commentLine = `${commentSymbol} File: ${relativePath}`;

        // 检查第一行是否已经是路径注释
        // 匹配模式：以注释符开头，包含 "File:" 字样
        const hasPathComment = lines[0] && lines[0].startsWith(commentSymbol) && lines[0].includes('File:');

        let newContent = '';

        if (hasPathComment) {
            // 情况 A: 已经有注释 -> 更新它 (防止文件移动后路径不对)
            if (lines[0] === commentLine) {
                // 内容一样，无需写入
                return; 
            }
            console.log(`[UPDATE] ${relativePath}`);
            lines[0] = commentLine; // 替换第一行
            newContent = lines.join('\n');
        } else {
            // 情况 B: 没有注释 -> 插入到第一行
            console.log(`[INSERT] ${relativePath}`);
            newContent = commentLine + '\n' + content;
        }

        fs.writeFileSync(filePath, newContent, 'utf8');

    } catch (err) {
        console.error(`Error processing ${filePath}:`, err.message);
    }
}

function walkDir(dir) {
    let entries;
    try {
        entries = fs.readdirSync(dir, { withFileTypes: true });
    } catch (err) {
        return;
    }

    for (const entry of entries) {
        if (shouldIgnore(entry.name)) continue;

        const fullPath = path.join(dir, entry.name);

        if (entry.isDirectory()) {
            walkDir(fullPath);
        } else if (entry.isFile()) {
            processFile(fullPath);
        }
    }
}

console.log('Adding relative path comments to files...\n');
walkDir(ROOT_DIR);
console.log('\nDone.');
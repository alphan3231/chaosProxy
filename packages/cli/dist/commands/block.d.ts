interface BlockOptions {
    redis: string;
}
export declare function block(ip: string, options: BlockOptions): Promise<void>;
export declare function unblock(ip: string, options: BlockOptions): Promise<void>;
export declare function listBlocked(options: BlockOptions): Promise<void>;
export {};
//# sourceMappingURL=block.d.ts.map
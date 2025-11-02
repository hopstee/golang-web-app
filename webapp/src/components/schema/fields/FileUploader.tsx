import { cn } from "@/lib/utils";
import { ImagePlus, Trash2 } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { Spinner } from "../../ui/spinner";
import { AspectRatio } from "../../ui/aspect-ratio";
import type { Field } from "@/types/pages";
import { Button } from "@/components/ui/button";

interface FileUploaderProps {
    field: Field;
    value: string;
    onChange: (name: string, value: string) => void;
}

export default function FileUploader(props: FileUploaderProps) {
    const [fileExists, setFileExists] = useState<boolean>(false);
    const [loadingImg, setLoadingImg] = useState<boolean>(false);
    const [isOverDropZone, setIsOverDropZone] = useState<boolean>(false);

    const inputRef = useRef<HTMLInputElement | null>(null);
    const imgRef = useRef<HTMLImageElement | null>(null);

    useEffect(() => {
        if (props.value && imgRef.current) {
            imgRef.current.src = props.value;
            setFileExists(!!props.value);
        }
    }, [props.value]);

    const handleRemoveImage = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        e.stopPropagation();
        const image = imgRef.current;
        if (image) {
            image.src = "";
            setFileExists(false);
            props.onChange(props.field.name, "");
        }
    }

    const handleDragEnter = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setIsOverDropZone(true);
    }

    const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setIsOverDropZone(false);
    }

    const handleClickDropZone = () => {
        const input = inputRef.current;
        if (input) {
            input.click();
        }
    }

    const handleDropFile = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();

        const file = e.dataTransfer.files[0];
        setImg(file);
    }

    const handleChangeFile = (e: React.ChangeEvent<HTMLInputElement>) => {
        e.preventDefault();
        e.stopPropagation();

        const file = e.target.files?.[0];
        setImg(file);
    }

    const setImg = (file?: File) => {
        if (!file) return;
        setLoadingImg(true);

        const image = imgRef.current
        if (image) {
            setFileExists(true);
            const src = URL.createObjectURL(file);
            image.src = src;
            props.onChange(props.field.name, src);
        }

        setLoadingImg(false);
    }

    return (
        <AspectRatio
            ratio={16 / 9}
            className={cn(
                "w-full dark:bg-input/30 shadow-xs rounded-md transition-all cursor-cell overflow-hidden",
                "border border-input hover:border-ring hover:ring-[3px] hover:ring-ring/50",
                "flex flex-col gap-2 items-center justify-center",
                isOverDropZone && "border-ring ring-[3px] ring-ring/50",
            )}
            onClick={handleClickDropZone}
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            onDrop={handleDropFile}
        >
            {fileExists && !loadingImg && (
                <Button
                    variant="destructive"
                    className="absolute top-2 right-2 z-50 p-1"
                    size="icon-sm"
                    onClick={handleRemoveImage}
                >
                    <Trash2 />
                </Button>
            )}

            <input ref={inputRef} type="file" accept="image/*" className="hidden h-full" onChange={handleChangeFile} />

            <img ref={imgRef} className={(fileExists && !loadingImg) ? "w-full h-full object-contain" : ""} />
            {loadingImg && (
                <div className="flex flex-col items-center gap-2 p-4">
                    <Spinner />
                    <p className="text-xs text-muted-foreground cursor-cell text-center">
                        Загружаем изображение...
                    </p>
                </div>
            )}
            {!fileExists && !loadingImg && (
                <div className="flex flex-col items-center gap-2 p-4">
                    <ImagePlus className="text-muted-foreground" />
                    <p className="text-xs text-muted-foreground cursor-cell text-center">
                        <span className="text-foreground">Перетащи</span> изображение сюда либо <span className="text-foreground">нажми</span> для выбора файла
                    </p>
                </div>
            )}
        </AspectRatio>
    );
}
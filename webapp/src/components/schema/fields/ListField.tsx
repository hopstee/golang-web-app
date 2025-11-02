import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import type { Field } from "@/types/pages";
import { Plus, Trash2 } from "lucide-react";
import FieldsScaffold from "../FieldsScaffold";
import { getDefaultValue } from "@/lib/utils";

interface ListFieldProps {
    type: "list[string]" | "list[text]" | "list[object]";
    field: Field;
    value: unknown[];
    handleChange: (name: string, value: unknown) => void;
    handleDelete: (name: string, index: number) => void;
}

export default function ListField(props: ListFieldProps) {
    switch (props.type) {
        case "list[string]":
            return renderStringsList(props);
        case "list[text]":
            return renderTextsList(props);
        case "list[object]":
            return renderObjectsList(props);
    }
}

function renderStringsList(props: ListFieldProps) {
    return (
        <div className="border-l pl-4 space-y-3">
            {props.value && (props.value as string[]).map((v, i) => (
                <div key={i} className="flex items-center gap-2">
                    <Input
                        value={v}
                        placeholder={props.field.label}
                        onChange={(e) => {
                            const updated = [...(props.value as string[])]
                            updated[i] = e.target.value
                            props.handleChange(props.field.name, updated)
                        }}
                    />
                    <Button
                        size="icon-sm"
                        variant="destructive"
                        onClick={() => props.handleDelete(props.field.name, i)}
                    >
                        <Trash2 />
                    </Button>
                </div>
            ))}
            <Button
                onClick={() => props.handleChange(props.field.name, [...(props.value as string[]), ""])}
                variant="secondary"
                size="sm"
            >
                <Plus />
                Добавить запись
            </Button>
        </div>
    );
}

function renderTextsList(props: ListFieldProps) {
    return (
        <div className="border-l pl-4 space-y-3">
            {props.value && (props.value as string[]).map((v, i) => (
                <div key={i} className="flex items-start gap-2">
                    <Textarea
                        key={i}
                        value={v}
                        onChange={(e) => {
                            const updated = [...(props.value as string[])]
                            updated[i] = e.target.value
                            props.handleChange(props.field.name, updated)
                        }}
                    />
                    <Button
                        size="icon-sm"
                        variant="destructive"
                        onClick={() => props.handleDelete(props.field.name, i)}
                    >
                        <Trash2 />
                    </Button>
                </div>
            ))}
            <Button
                onClick={() => props.handleChange(props.field.name, [...(props.value as string[]), ""])}
                variant="secondary"
                size="sm"
            >
                <Plus />
                Добавить текст
            </Button>
        </div>
    );
}

function renderObjectsList(props: ListFieldProps) {
    return (
        <div className="border-l pl-4 space-y-3">
            {props.value && (props.value as Record<string, unknown>[]).map((item, index) => (
                <div key={index} className="border-b space-y-2 pb-6">
                    <div className="flex items-center justify-between">
                        <Label>
                            {props.field.label} #{index + 1}
                        </Label>
                        <Button
                            type="button"
                            variant="destructive"
                            size="sm"
                            onClick={() => props.handleDelete(props.field.name, index)}
                        >
                            <Trash2 />
                            Удалить
                        </Button>
                    </div>
                    <FieldsScaffold
                        fields={props.field.schema?.fields || []}
                        data={item}
                        onChange={(updatedItem) => {
                            const updated = [...(props.value as Record<string, unknown>[])];
                            updated[index] = updatedItem;
                            props.handleChange(props.field.name, updated);
                        }}
                    />
                </div>
            ))}
            <Button
                onClick={() =>
                    props.handleChange(props.field.name, [
                        ...(props.value as Record<string, unknown>[]),
                        props.field.schema?.fields.reduce(
                            (acc, f) => ({ ...acc, [f.name]: getDefaultValue(f.type) }),
                            {}
                        ) ?? {},
                    ])
                }
                variant="secondary"
                size="sm"
            >
                <Plus />
                Добавить объект
            </Button>
        </div>
    );
}

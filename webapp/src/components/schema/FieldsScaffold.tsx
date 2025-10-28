import type { Field } from "@/types/pages";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "../ui/textarea";
import { useEffect, useState } from "react";
import { Checkbox } from "../ui/checkbox";
import { Button } from "../ui/button";
import { Plus, Trash2 } from "lucide-react";

interface FieldsScaffoldProps {
    fields: Field[];
    data: Record<string, unknown>
    onChange?: (updated: Record<string, unknown>) => void
}

export default function FieldsScaffold(props: FieldsScaffoldProps) {
    const {
        fields,
        data,
        onChange
    } = props;

    const [values, setValues] = useState<Record<string, unknown>>(() =>
        fields.reduce((acc, f) => ({ ...acc, [f.name]: data[f.name] ?? getDefaultValue(f.type) }), {})
    )

    useEffect(() => {
        setValues(
            fields.reduce(
                (acc, f) => ({ ...acc, [f.name]: data[f.name] ?? getDefaultValue(f.type) }),
                {}
            )
        )
    }, [data, fields])

    const handleChange = (name: string, value: unknown) => {
        const updated = { ...values, [name]: value }
        setValues(updated);
        onChange?.(updated);
    }

    const handleDelete = (name: string, index: number) => {
        const list = [...((values[name] as unknown[]) || [])]
        list.splice(index, 1);
        handleChange(name, list);
    }

    const renderField = (field: Field) => {
        const value = values[field.name] ?? ""

        switch (field.type) {
            case "string":
                return (
                    <Input
                        value={value as string}
                        onChange={(e) => handleChange(field.name, e.target.value)}
                        placeholder={field.label}
                    />
                );
            case "text":
                return (
                    <Textarea
                        value={value as string}
                        onChange={(e) => handleChange(field.name, e.target.value)}
                        placeholder={field.label}
                    />
                );
            case "bool":
            case "boolean":
                return (
                    <Checkbox
                        checked={value === "true"}
                        onCheckedChange={(checked) => handleChange(field.name, String(checked))}
                    />
                );
            case "list[string]":
                return (
                    <div className="border-l pl-4 space-y-3">
                        {value && (value as string[]).map((v, i) => (
                            <div key={i} className="flex items-center gap-2">
                                <Input
                                    value={v}
                                    placeholder={field.label}
                                    onChange={(e) => {
                                        const updated = [...(value as string[])]
                                        updated[i] = e.target.value
                                        handleChange(field.name, updated)
                                    }}
                                />
                                <Button
                                    size="icon-sm"
                                    variant="destructive"
                                    onClick={() => handleDelete(field.name, i)}
                                >
                                    <Trash2 />
                                </Button>
                            </div>
                        ))}
                        <Button
                            onClick={() => handleChange(field.name, [...(value as string[]), ""])}
                            variant="secondary"
                            size="sm"
                        >
                            <Plus />
                            Добавить запись
                        </Button>
                    </div>
                );
            case "list[text]":
                return (
                    <div className="border-l pl-4 space-y-3">
                        {value && (value as string[]).map((v, i) => (
                            <div key={i} className="flex items-center gap-2">
                                <Textarea
                                    key={i}
                                    value={v}
                                    onChange={(e) => {
                                        const updated = [...(value as string[])]
                                        updated[i] = e.target.value
                                        handleChange(field.name, updated)
                                    }}
                                />
                                <Button
                                    size="icon-sm"
                                    variant="destructive"
                                    onClick={() => handleDelete(field.name, i)}
                                >
                                    <Trash2 />
                                </Button>
                            </div>
                        ))}
                        <Button
                            onClick={() => handleChange(field.name, [...(value as string[]), ""])}
                            variant="secondary"
                            size="sm"
                        >
                            <Plus />
                            Добавить текст
                        </Button>
                    </div>
                );

            case "list[object]":
                return (
                    <div className="border-l pl-4 space-y-3">
                        {value && (value as Record<string, unknown>[]).map((item, index) => (
                            <div key={index} className="border-b space-y-2 pb-6">
                                <div className="flex items-center justify-between">
                                    <Label>
                                        {field.label} #{index + 1}
                                    </Label>
                                    <Button
                                        type="button"
                                        variant="destructive"
                                        size="sm"
                                        onClick={() => handleDelete(field.name, index)}
                                    >
                                        <Trash2 />
                                        Удалить
                                    </Button>
                                </div>
                                <FieldsScaffold
                                    fields={field.schema?.fields || []}
                                    data={item}
                                    onChange={(updatedItem) => {
                                        const updated = [...(value as Record<string, unknown>[])];
                                        updated[index] = updatedItem;
                                        handleChange(field.name, updated);
                                    }}
                                />
                            </div>
                        ))}
                        <Button
                            onClick={() =>
                                handleChange(field.name, [
                                    ...(value as Record<string, unknown>[]),
                                    field.schema?.fields.reduce(
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
    }

    if (fields.length === 0) {
        return <div>Нет доступных полей для редактирования</div>
    }

    return (
        <div className="space-y-6">
            {fields.map(field => (
                <div className="space-y-2" key={field.id}>
                    <Label key={field.id}>{field.label}</Label>
                    {renderField(field)}
                </div>
            ))}
        </div>
    )
}

function getDefaultValue(type: string): unknown {
    if (type.startsWith("list[")) return [];
    if (type === "bool" || type === "boolean") return false;
    return "";
}

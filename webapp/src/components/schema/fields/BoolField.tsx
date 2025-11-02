import { Checkbox } from "@/components/ui/checkbox";
import type { Field } from "@/types/pages";

interface BoolFieldProps {
    field: Field;
    value: string;
    onChange: (name: string, value: string) => void;
}

export default function BoolField(props: BoolFieldProps) {
    return (
        <Checkbox
            checked={props.value === "true"}
            onCheckedChange={(checked) => props.onChange(props.field.name, String(checked))}
        />
    );
}
import { Button } from "@/components/ui/button";
import { ButtonGroup } from "@/components/ui/button-group";
import { Input } from "@/components/ui/input";
import type { Field } from "@/types/entities";
import { Minus, Plus } from "lucide-react";

interface NumberFieldProps {
    field: Field;
    value: number;
    onChange: (name: string, value: number) => void;
}

const maxValue = 999999999999;
const minValue = 0;

export default function NumberField(props: NumberFieldProps) {
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = Number(e.target.value);
        if (isNaN(value)) return;
        if (value > maxValue) {
            props.onChange(props.field.name, maxValue);
            return;
        }
        if (value < minValue) {
            props.onChange(props.field.name, minValue);
            return;
        }
        props.onChange(props.field.name, value)
    }

    const handleIncrement = () => {
        if (props.value === maxValue) return;
        props.onChange(props.field.name, props.value + 1);
    }

    const handleDecrement = () => {
        if (props.value === minValue) return;
        props.onChange(props.field.name, props.value - 1);
    }

    return (
        <ButtonGroup className="w-full">
            <Input
                pattern="[0-9]*"
                value={props.value}
                onChange={handleChange}
                placeholder={props.field.label}
            />
            <Button size="icon" variant="outline" onClick={handleDecrement} disabled={props.value <= minValue}><Minus /></Button>
            <Button size="icon" variant="outline" onClick={handleIncrement} disabled={props.value >= maxValue}><Plus /></Button>
        </ButtonGroup>
    );
}